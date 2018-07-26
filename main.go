package main

import (
	"github.com/tonyStreet/projectOrder/db"
	"github.com/tonyStreet/projectOrder/server"
	"log"
	"github.com/tonyStreet/projectOrder/config"
)

func main() {
	confFile := "config.yml"
	config.InitConfig(confFile)
	err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(server.Dispatch())
}
