package routes

import (
	"github.com/aakritigkmit/payment-gateway/internal/handlers"
	"github.com/aakritigkmit/payment-gateway/internal/repository"
	"github.com/aakritigkmit/payment-gateway/internal/services"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupAuthRoutes registers authentication-related routes
func SetupAuthRoutes(r chi.Router, db *mongo.Database) {
	userRepo := repository.NewUserRepo(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)
}
