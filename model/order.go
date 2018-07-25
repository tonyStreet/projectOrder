package model

import (
	"errors"
	"github.com/tonyStreet/projectOrder/db"
	"strings"
)

const (
	ERROR_MISSING_ORIGIN       = "Missing origin"
	ERROR_ORIGIN_TYPE          = "origin must be of type string array"
	ERROR_ORIGIN_VALUE         = "origin must contain latitude and longitude"
	ERROR_MISSING_DESTINATION  = "Missing destination"
	ERROR_DESTINATION_TYPE     = "destination must be of type string array"
	ERROR_DESTINATION_VALUE    = "destination must contain latitude and longitude"
	ORDER_ALREADY_BEEN_TAKEN   = "ORDER_ALREADY_BEEN_TAKEN"
	ORDER_STATUS_UNASSIGNED    = "UNASSIGN"
	ORDER_STATUS_TAKEN         = "TAKEN"
	INVALID_ORDER_STATUS       = "Invalid order status"
	ERROR_ORDER_STATUS_TYPE    = "status must be of type string"
	ERROR_MISSING_ORDER_STATUS = "Missing required parameter 'status'"
	INVALID_JSON_FORMAT        = "Invalid json format"
	ERROR_ORDER_DB_SAVE        = "Error on saving order"
	ORDER_NOT_EXISTS           = "Order doesn't exists"
)

type Order struct {
	ID       int64  `json:"id"`
	Distance int64  `json:"distance"`
	Status   string `json:"status"`
}

type CreateOrderRequest struct {
	Origin      []string `json:"origin"`
	Destination []string `json:"destination"`
}

type CreateOrderResponse struct {
	Origin      []string `json:"origin"`
	Destination []string `json:"destination"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (order *Order) AddOrder() error {
	datasource, err := db.GetDataSource()
	if err != nil {
		return err
	}
	res, err := datasource.Exec(`INSERT INTO orders (status,distance) VALUES (?, ?)`, strings.ToLower(ORDER_STATUS_UNASSIGNED), order.Distance)
	if err != nil {
		return err
	} else {
		id, err := res.LastInsertId()
		if err != nil {
			return err
		} else {
			order.ID = id
		}
	}
	return err
}

func (order *Order) GetOrderStatus() (status string, err error) {
	datasource, err := db.GetDataSource()
	if err != nil {
		return status, err
	}
	err = datasource.QueryRow("SELECT status FROM orders where id = ?", order.ID).Scan(&status)
	return status, err
}

func (order *Order) UpdateOrderStatus() error {
	datasource, err := db.GetDataSource()
	if err != nil {
		return err
	}
	res, err := datasource.Exec("UPDATE orders SET status = ? where id = ?", order.Status, order.ID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected > int64(1) {
		return errors.New("More than one order updated")
	}
	return err
}

func GetOrderList(page int, limit int) (orders []Order, err error) {
	offset := (page - 1) * limit
	datasource, err := db.GetDataSource()
	if err != nil {
		return orders, err
	}
	rows, err := datasource.Query("SELECT id, distance, status FROM orders ORDER BY ID ASC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return orders, err
	}
	for rows.Next() {
		var id int64
		var distance int64
		var status string
		rows.Scan(&id, &distance, &status)
		orders = append(orders, Order{id, distance, status})
	}
	return orders, err
}
