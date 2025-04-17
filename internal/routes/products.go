package routes

import (
	"github.com/aakritigkmit/payment-gateway/internal/handlers"
	middlewares "github.com/aakritigkmit/payment-gateway/internal/middleware"
	"github.com/aakritigkmit/payment-gateway/internal/repository"
	"github.com/aakritigkmit/payment-gateway/internal/services"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupProductRoutes(r chi.Router, db *mongo.Database) {
	productRepo := repository.NewProductRepo(db)
	productTransactionRepo := repository.NewProductTransactionRepo(db)
	productOrderRepo := repository.NewProductOrderRepo(db)

	productService := services.NewProductService(productRepo, productTransactionRepo, productOrderRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Define routes
	r.With(middlewares.AuthMiddleware).Post("/sync", productHandler.SyncProducts)
	r.With(middlewares.AuthMiddleware).Post("/transaction", productHandler.HandleProductTransaction)
	r.With(middlewares.AuthMiddleware).Post("/transactions/bulk", productHandler.CreateBulkProductTransaction)
}
