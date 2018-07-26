package config

import "testing"

func TestInitConfig(t *testing.T) {
	confFile := "../config.yml"
	err := InitConfig(confFile)
	if err != nil {
		t.Error(err)
	}
}

func TestGetConfig(t *testing.T) {
	confFile := "../config.yml"
	err := InitConfig(confFile)
	if err != nil {
		t.Error(err)
	}
	conf := GetConfig()
	if conf.DB.Name != "logistics" {
		t.Fail()
	}
}
