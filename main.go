package main

import (
	"github.com/tonyStreet/projectOrder/server"
	"log"
)

func main() {
	log.Fatal(server.Dispatch())
}
