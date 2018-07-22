package model

const (
	ERROR_ORIGIN_TYPE = "origin must be of type string array"
	ERROR_ORIGIN_VALUE = "origin must contain latitude and longitude"
	ERROR_MISSING_DESTINATION = "Missing destination"
	ERROR_DESTINATION_TYPE = "destination must be of type string array"
	ERROR_DESTINATION_VALUE = "destination must contain latitude and longitude"
)

type CreateOrderRequest struct {
	Origin      []string `json:"origin"`
	Destination []string `json:"destination"`
}

type CreateOrderSuccessResponse struct {
	Error string `json:"error"`
}

type CreateOrderErrorResponse struct {
	Error string `json:"error"`
}
