package main

import (
	"log"
	"net/http"

	"github.com/smjn/ipl18/backend"
	"github.com/smjn/ipl18/scraper"
)

func main() {
	kill := make(chan bool)
	go scraper.Start(kill)
	log.Println("Starting server on port 3000...")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", backend.SetupAndGetRouter()))
}
