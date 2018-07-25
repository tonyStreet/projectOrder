package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/tonyStreet/projectOrder/model"
	"github.com/tonyStreet/projectOrder/service"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	reqbody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		errMsg := model.ErrorResponse{err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	var req map[string]interface{}

	if err = json.Unmarshal(reqbody, &req); err != nil {
		errMsg := model.ErrorResponse{model.INVALID_JSON_FORMAT}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	var origin []string
	var org []interface{}
	if _, ok := req["origin"]; ok {
		org, ok = req["origin"].([]interface{})
		if !ok {
			errMsg := model.ErrorResponse{model.ERROR_ORIGIN_TYPE}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errMsg)
			return
		}
		if len(org) < 2 {
			errMsg := model.ErrorResponse{model.ERROR_ORIGIN_VALUE}
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
			errMsg := model.ErrorResponse{model.ERROR_DESTINATION_TYPE}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errMsg)
			return
		}
		if len(des) < 2 {
			errMsg := model.ErrorResponse{model.ERROR_DESTINATION_VALUE}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errMsg)
			return
		}
		for _, v := range des {
			destination = append(destination, v.(string))
		}
	}

	if len(origin) == 0 {
		errMsg := model.ErrorResponse{model.ERROR_MISSING_ORIGIN}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	if len(destination) == 0 {
		errMsg := model.ErrorResponse{model.ERROR_MISSING_DESTINATION}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	originStr := strings.TrimSpace(strings.Join(origin, ","))
	destinationStr := strings.TrimSpace(strings.Join(destination, ","))

	distance, err := service.GetDistance(originStr, destinationStr)

	if err != nil {
		errMsg := model.ErrorResponse{err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	order := &model.Order{Distance: distance, Status: model.ORDER_STATUS_UNASSIGNED}
	err = order.AddOrder()
	if err != nil {
		errMsg := model.ErrorResponse{model.ERROR_ORDER_DB_SAVE}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
	return

}

func TakeOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	reqbody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		errMsg := model.ErrorResponse{err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	var req map[string]interface{}

	if err = json.Unmarshal(reqbody, &req); err != nil {
		errMsg := model.ErrorResponse{model.INVALID_JSON_FORMAT}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	status := ""
	if _, ok := req["status"]; ok {
		status, ok = req["status"].(string)
		if !ok {
			errMsg := model.ErrorResponse{model.ERROR_ORDER_STATUS_TYPE}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errMsg)
			return
		}
	}

	if status == "" {
		errMsg := model.ErrorResponse{model.ERROR_MISSING_ORDER_STATUS}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	if strings.ToLower(status) != strings.ToLower(model.ORDER_STATUS_TAKEN) {
		errMsg := model.ErrorResponse{model.INVALID_ORDER_STATUS}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	orderID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		errMsg := model.ErrorResponse{err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	order := &model.Order{ID: orderID}
	orderStatus, err := order.GetOrderStatus()
	if err != nil {
		errMsg := model.ErrorResponse{model.ORDER_NOT_EXISTS}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	orderStatus = strings.ToUpper(orderStatus)
	if orderStatus == model.ORDER_STATUS_TAKEN {
		errMsg := model.ErrorResponse{model.ORDER_ALREADY_BEEN_TAKEN}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	order.Status = status
	err = order.UpdateOrderStatus()

	if err != nil {
		errMsg := model.ErrorResponse{err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{status:"success"}`))
}

func ListOrderHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.FormValue("page")
	limitStr := r.FormValue("limit")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errMsg := model.ErrorResponse{err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		errMsg := model.ErrorResponse{err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	orders, err := model.GetOrderList(page, limit)
	if err != nil {
		errMsg := model.ErrorResponse{err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
