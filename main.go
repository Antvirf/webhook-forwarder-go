package main

import (
	"log"

	"github.com/Antvirf/webhook-forwarder-go/api"
)

const (
	serverAddress = "0.0.0.0:8080"
)

func main() {
	server := api.NewServer()

	err := server.Start(serverAddress)

	if err != nil {
		log.Fatal("cannot start server")
	}

	// r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	// r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
