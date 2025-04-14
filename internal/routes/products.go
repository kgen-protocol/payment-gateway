package routes

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupProductRoutes(r chi.Router, db *mongo.Database) {
	// orderRepo := repository.NewOrderRepo(db)
	productsRepo := repository.NewProductRepo(db)
	// orderService := services.NewOrderService(orderRepo, transactionRepo)
	// orderHandler := handlers.NewOrderHandler(orderService)

	// r.Use(middlewares.AuthMiddleware) // Apply auth middleware

	r.Post("/place", productHandler.performTransactionHandler)
	r.Post("/callback/order-status", productHandler.confirmTransactionHandler)
	// r.get("/:id" , orderHandler.)
}
