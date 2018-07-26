package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var config *Config

type Config struct {
	GoogleMapsAPIKey string `yaml:"google_maps_api_key"`
	DB               struct {
		IP       string `yaml:"ip"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	}
}

func InitConfig(configFile string) error {
	payloadByte, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	payload := string(payloadByte)
	err = yaml.Unmarshal([]byte(payload), &config)
	if err != nil {
		return err
	}
	return err
}

func GetConfig() *Config {
	return config
}
