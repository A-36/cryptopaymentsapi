package main

import (
	"cryptoapi/core/config"
	"cryptoapi/core/masters/webserver"
	"log"
	"os"
)

func init() {
	if err := config.Load("assets/config.json"); err != nil {
		log.Printf("Could not parse config.json.\r\nYou have made a mistake in the JSON format.")
		os.Exit(1)
	}
	go config.WatchChanges("assets/config.json")
}

func main() {
	log.Printf("Crypto Payments API Status: Online")
	webserver.Start()
}
