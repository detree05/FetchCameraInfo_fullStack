package main

import (
	"fmt"
	"log"
	"net/http"

	"fci-backend.detree05.com/cfg"
)

var LogLine uint = 1

func main() {
	log.Printf("[~] Reading configuration file...")

	err := cfg.ReadConfigurationFile()
	if err != nil {
		log.Fatalf("[!] Error reading configuration file! %s", err)
	}

	port := fmt.Sprintf(":%d", cfg.Config.Server.Port)
	log.Printf("[~] Starting server at port %s", port)

	http.HandleFunc("/", invalidReq)
	http.HandleFunc("/ping", healthCheck)
	http.HandleFunc("/getChannel", getChannel)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
