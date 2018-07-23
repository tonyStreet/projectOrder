package handler

import (
	"bytes"
	"encoding/json"
	"github.com/tonyStreet/projectOrder/model"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
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
	ResponseBody string
}

func TestCreateOrder(t *testing.T) {
	// Create a request to pass to our handler.
	url := "/order"
	requests := []TestRequest{
		TestRequest{Request{"POST", `{"origin": "41.43206"}`}, Response{400, model.ERROR_ORIGIN_TYPE}},
		TestRequest{Request{"POST", `{"origin": ["41.43206"]}`}, Response{400, model.ERROR_ORIGIN_VALUE}},
		TestRequest{Request{"POST", `{"origin": ["41.43206","-81.38992"]}`}, Response{400, model.ERROR_MISSING_DESTINATION}},
		TestRequest{Request{"POST", `{"destination": "40.714224"}`}, Response{400, model.ERROR_DESTINATION_TYPE}},
		TestRequest{Request{"POST", `{"destination": ["40.714224"]}`}, Response{400, model.ERROR_DESTINATION_VALUE}},
		TestRequest{Request{"POST", `{"destination": ["40.714224","-73.961452"]}`}, Response{400, model.ERROR_MISSING_ORIGIN}},
	}

	for _, r := range requests {
		var jsonStr = []byte(r.Req.RequestBody)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
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
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check the response body is what we expect.
		expected := r.Res.ResponseBody
		var responseBody map[string]interface{}
		if rr.Code == http.StatusBadRequest {
			json.Unmarshal([]byte(rr.Body.String()), &responseBody)
			response := responseBody["error"].(string)
			if response != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					response, expected)
			}
		} else if rr.Code == http.StatusOK {
			reqbody, err := ioutil.ReadAll(rr.Body)
			var req map[string]interface{}
			if err = json.Unmarshal(reqbody, &req); err != nil {
				errMsg := model.CreateOrderErrorResponse{err.Error()}
				t.Error(errMsg)
			}
		}
	}
}
