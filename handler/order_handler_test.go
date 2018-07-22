package handler

import (
	"bytes"
	"encoding/json"
	"github.com/tonyStreet/projectOrder/model"
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
		TestRequest{Request{"POST", `{"origin": "14.5995"}`}, Response{400, model.ERROR_ORIGIN_TYPE}},
		TestRequest{Request{"POST",`{"origin": ["14.5995"]}`}, Response{400, model.ERROR_ORIGIN_VALUE}},
		TestRequest{Request{"POST",`{"origin": ["14.5995","120.9842"]}`}, Response{400, model.ERROR_MISSING_DESTINATION}},
		TestRequest{Request{"POST", `{"destination": "22.3390408802915"}`}, Response{400, model.ERROR_DESTINATION_TYPE}},
		TestRequest{Request{"POST", `{"destination": ["22.3390408802915"]}`}, Response{400, model.ERROR_DESTINATION_VALUE}},
	}

	/*`{"origin":["44.968046"]}`,
	`{"origin":["44.968046", "-94.420307"]}`,
	`{"destination":["22.3390408802915", "114.1486719802915"]}`,*/

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
		}
	}
}
