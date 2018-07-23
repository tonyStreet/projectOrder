package handler

import (
	"encoding/json"
	"github.com/tonyStreet/projectOrder/model"
	"io/ioutil"
	"net/http"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	reqbody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		errMsg := model.CreateOrderErrorResponse{err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	var req map[string]interface{}

	if err = json.Unmarshal(reqbody, &req); err != nil {
		errMsg := model.CreateOrderErrorResponse{err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	var origin []string
	var org []interface{}
	if _, ok := req["origin"]; ok {
		org, ok = req["origin"].([]interface{})
		if !ok {
			errMsg := model.CreateOrderErrorResponse{model.ERROR_ORIGIN_TYPE}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errMsg)
			return
		}
		if len(org) < 2 {
			errMsg := model.CreateOrderErrorResponse{model.ERROR_ORIGIN_VALUE}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errMsg)
			return
		}
		for _, v := range org {
			origin = append(origin, v.(string))
		}
	}
	var destination []string
	var des []interface{}
	if _, ok := req["destination"]; ok {
		des, ok = req["destination"].([]interface{})
		if !ok {
			errMsg := model.CreateOrderErrorResponse{model.ERROR_DESTINATION_TYPE}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errMsg)
			return
		}
		if len(des) < 2 {
			errMsg := model.CreateOrderErrorResponse{model.ERROR_DESTINATION_VALUE}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errMsg)
			return
		}
		for _, v := range des {
			destination = append(destination, v.(string))
		}
	}

	if len(origin) == 0 {
		errMsg := model.CreateOrderErrorResponse{model.ERROR_MISSING_ORIGIN}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	if len(destination) == 0 {
		errMsg := model.CreateOrderErrorResponse{model.ERROR_MISSING_DESTINATION}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errMsg)
		return
	}
}
