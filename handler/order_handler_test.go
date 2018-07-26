package handler

import (
	"bytes"
	"encoding/json"
	"github.com/tonyStreet/projectOrder/db"
	"github.com/tonyStreet/projectOrder/model"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"github.com/tonyStreet/projectOrder/config"
	"github.com/gorilla/mux"
)

type TestRequest struct {
	Req Request
	Res Response
}

type Request struct {
	Method      string
	RequestBody string
}

type Response struct {
	Code         int
	ResponseBody interface{}
}

//Test Requires connection to db
func TestCreateOrder(t *testing.T) {
	// Create a request to pass to our handler.
	confFile := "../config_test.yml"
	config.InitConfig(confFile)
	db.InitDB()
	url := "/order"
	requests := []TestRequest{
		TestRequest{Request{"POST", `{"origin": "41.43206"}`}, Response{http.StatusBadRequest, model.ERROR_ORIGIN_TYPE}},
		TestRequest{Request{"POST", `{"origin": ["41.43206"]}`}, Response{http.StatusBadRequest, model.ERROR_ORIGIN_VALUE}},
		TestRequest{Request{"POST", `{"origin": ["41.43206","-81.38992"]}`}, Response{http.StatusBadRequest, model.ERROR_MISSING_DESTINATION}},
		TestRequest{Request{"POST", `{"destination": "40.714224"}`}, Response{http.StatusBadRequest, model.ERROR_DESTINATION_TYPE}},
		TestRequest{Request{"POST", `{"destination": ["40.714224"]}`}, Response{http.StatusBadRequest, model.ERROR_DESTINATION_VALUE}},
		TestRequest{Request{"POST", `{"destination": ["40.714224","-73.961452"]}`}, Response{http.StatusBadRequest, model.ERROR_MISSING_ORIGIN}},
		TestRequest{Request{"POST", `{"origin" : ["41.43206","-81.38992"], "destination": ["40.714224","-73.961452"]}`}, Response{http.StatusOK, model.Order{0, 710535, model.ORDER_STATUS_UNASSIGNED}}},
	}

	for testNum, r := range requests {
		var jsonStr = []byte(r.Req.RequestBody)
		req, err := http.NewRequest(r.Req.Method, url, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateOrder)

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != r.Res.Code {
			t.Errorf("Test num: %v : create order  handler returned wrong status code: got %v want %v",
				testNum, status, http.StatusOK)
		}

		// Check the response body is what we expect.
		expected := r.Res.ResponseBody
		var responseBody map[string]interface{}
		if rr.Code == http.StatusBadRequest {
			json.Unmarshal([]byte(rr.Body.String()), &responseBody)
			response := responseBody["error"].(string)
			if response != expected.(string) {
				t.Errorf("Test num: %v : create order handler returned unexpected body: got %v want %v",
					testNum, response, expected)
			}
		} else if rr.Code == http.StatusOK {
			expectedResponse := expected.(model.Order)
			res := model.Order{}
			err := json.Unmarshal([]byte(rr.Body.String()), &res)
			if err != nil {
				t.Error(err.Error())
			}
			if res.ID == 0 {
				t.Error("Create order request not saved in db")
			}
			expectedResponse.ID = res.ID
			if !reflect.DeepEqual(res, expectedResponse) {
				t.Fail()
			}
		}
	}
}


func TestTakeOrderHandler(t *testing.T) {
	// Create a request to pass to our handler.
	confFile := "../config_test.yml"
	config.InitConfig(confFile)
	db.InitDB()
	url := "/order/1"
	requests := []TestRequest{
		TestRequest{Request{"PUT", `{"status":"taken"}`}, Response{http.StatusOK, "success"}},
	}

	for testNum, r := range requests {
		var jsonStr = []byte(r.Req.RequestBody)
		req, err := http.NewRequest(r.Req.Method, url, bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()

		router := mux.NewRouter()

		router.HandleFunc("/order/{id:[0-9]+}", TakeOrderHandler).Methods(http.MethodPut)

		router.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != r.Res.Code {
			t.Errorf("Test num: %v : take order  handler returned wrong status code: got %v want %v",
				testNum, status, http.StatusOK)
		}
		// Check the response body is what we expect.
		expected := r.Res.ResponseBody
		var responseBody map[string]interface{}
		if rr.Code == http.StatusOK {
			json.Unmarshal([]byte(rr.Body.String()), &responseBody)
			response := responseBody["status"].(string)
			if response != expected {
				t.Fail()
			}
		}
	}
}

func TestListOrderHandler(t *testing.T) {
	// Create a request to pass to our handler.
	confFile := "../config_test.yml"
	config.InitConfig(confFile)
	db.InitDB()
	url := "/orders?page=1&limit=10"
	requests := []TestRequest{
		TestRequest{Request{"GET", ``}, Response{http.StatusOK, ""}},
	}

	for testNum, r := range requests {
		req, err := http.NewRequest(r.Req.Method, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/orders", ListOrderHandler).Queries("page", "{page:[0-9]+}", "limit", "{limit:[0-9]+}").Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != r.Res.Code {
			t.Errorf("Test num: %v : take order  handler returned wrong status code: got %v want %v",
				testNum, status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var responseBody []model.Order
		if rr.Code == http.StatusOK {
			json.Unmarshal([]byte(rr.Body.String()), &responseBody)
			if len(responseBody) > 10 {
				t.Fail()
			}
			if len(responseBody) > 0 {
				if responseBody[0].ID != 1 {
					t.Fail()
				}
			}
		}
	}
}
