package main

import (
	"log"
	"os"

	"github.com/Antvirf/webhook-forwarder-go/api"
)

const (
	serverAddress = "0.0.0.0:8000"
)

func main() {
	// Load environment variables
	os.Setenv("TARGET_URL", "http://localhost:8000/receive_webhook")
	os.Setenv("WEBHOOK_TOKEN_SECRET", "secret")

	TARGET_URL := os.Getenv("TARGET_URL")
	WEBHOOK_TOKEN_SECRET := os.Getenv("WEBHOOK_TOKEN_SECRET")

	if TARGET_URL == "" {
		log.Fatal("failed to start, empty target url")
	}
	if WEBHOOK_TOKEN_SECRET == "" {
		log.Println("warning: provided token is blank")
	}

	// Start server
	server := api.NewServer()
	err := server.Start(serverAddress)

	if err != nil {
		log.Fatal("cannot start server")
	}
}
