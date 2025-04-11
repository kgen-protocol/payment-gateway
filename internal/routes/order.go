package routes

import (
	"github.com/aakritigkmit/payment-gateway/internal/handlers"
	"github.com/aakritigkmit/payment-gateway/internal/repository"
	"github.com/aakritigkmit/payment-gateway/internal/services"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupOrderRoutes(r chi.Router, db *mongo.Database) {
	orderRepo := repository.NewOrderRepo(db)
	transactionRepo := repository.NewTransactionRepo(db)
	orderService := services.NewOrderService(orderRepo, transactionRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	// r.Use(middlewares.AuthMiddleware) // Apply auth middleware

	r.Post("/place", orderHandler.PlaceOrder)
	r.Post("/callback/order-status", orderHandler.HandleCallback)
}
