package db

import "testing"

//DB integration test
func TestDBConnection(t *testing.T) {
	err := InitDB()
	if err != nil {
		t.Error(err)
	}
	_, err = GetDataSource()
	if err != nil {
		t.Error(err)
	}
}
