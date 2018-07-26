package service

import (
	"testing"
	"github.com/tonyStreet/projectOrder/config"
)


//Test Requires internet connection to connect to google api
func TestGetDistance(t *testing.T) {
	confFile := "../config_test.yml"
	config.InitConfig(confFile)
	distance, err := GetDistance("41.43206,-81.38992", "40.714224,-73.961452")
	if err != nil {
		t.Fail()
	}
	if distance != 710535 {
		t.Fail()
	}
}
