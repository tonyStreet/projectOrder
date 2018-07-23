package service

import (
	"encoding/json"
	"github.com/tonyStreet/projectOrder/model"
	"net"
	"net/http"
	"time"
)

const (
	DIRECTION_API        = "https://maps.googleapis.com/maps/api/directions/"
	JSON_RESPONSE_FORMAT = "json"
	MAPS_API_KEY         = "AIzaSyCL2LEIXOFKsF_8zYg6SO48wSxpQb1StRg"
)

func GetDistance(origin string, destination string) (distance string, err error) {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	url := DIRECTION_API + JSON_RESPONSE_FORMAT + "?origin=" + origin + "&destination=" + destination + "&key=" + MAPS_API_KEY
	response, err := netClient.Get(url)
	if err != nil {
		return distance, err
	}
	var direction model.DirectionResponse
	err = json.NewDecoder(response.Body).Decode(&direction)
	if err != nil {
		return distance, err
	}
	if direction.Status == "OK" && len(direction.Routes) > 0 {
		distance = direction.Routes[0].Legs[0].Distance.Text
	}
	return distance, err
}
