package routes

import (
	"github.com/aakritigkmit/payment-gateway/internal/handlers"
	"github.com/aakritigkmit/payment-gateway/internal/repository"
	"github.com/aakritigkmit/payment-gateway/internal/services"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupProductRoutes(r chi.Router, db *mongo.Database) {
	productRepo := repository.NewProductRepo(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	r.Post("/sync", productHandler.FetchAndSaveProducts)
}
