package service

import "testing"

func TestGetDistance(t *testing.T) {
	distance, err := GetDistance("41.43206,-81.38992", "40.714224,-73.961452")
	if err != nil {
		t.Fail()
	}
	if distance != 710535 {
		t.Fail()
	}
}
