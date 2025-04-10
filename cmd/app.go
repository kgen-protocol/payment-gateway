package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aakritigkmit/payment-gateway/internal/config"
	"github.com/aakritigkmit/payment-gateway/internal/routes"
	"github.com/aakritigkmit/payment-gateway/internal/utils"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	router *chi.Mux
	db     *mongo.Database
}

// NewApp initializes a new application instance.
func NewApp() (*App, error) {
	// Initialize configuration.
	config.Init()

	// Connect to the database.
	db, err := utils.ConnectDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize Redis
	utils.InitRedis()
	if err := utils.PingRedis(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// Initialize router and setup routes.
	router := chi.NewRouter()
	routes.SetupRoutes(router, db)

	return &App{
		router: router,
		db:     db,
	}, nil
}

// Start runs the HTTP server with graceful shutdown.
func (a *App) Start(ctx context.Context) error {
	// Get port from our configuration singleton.
	cfg := config.GetConfig()
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: a.router,
	}

	// Log server start.
	log.Println("Server running on port", cfg.Port)

	// Start the server in a goroutine.
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.ListenAndServe()
	}()

	// Gracefully shutdown when context is done.
	select {
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		log.Println("Shutting down server...")
		return server.Shutdown(context.Background())
	}
}
