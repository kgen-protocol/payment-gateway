package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize the app
	app, err := NewApp()
	if err != nil {
		log.Fatalf("Error initializing application: %v", err)
	}

	// Handle graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start the application
	if err := app.Start(ctx); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
