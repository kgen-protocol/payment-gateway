package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mongoClient, router := initializeApp()
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatal("Error disconnecting MongoDB:", err)
		}
		fmt.Println("Disconnected from MongoDB")
	}()

	port := ":4444"
	fmt.Println("Server started on port", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
