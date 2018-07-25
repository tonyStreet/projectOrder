package main

import (
	"github.com/tonyStreet/projectOrder/db"
	"github.com/tonyStreet/projectOrder/server"
	"log"
)

func main() {
	err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(server.Dispatch())
}
