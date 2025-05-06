package routes

import (
	"github.com/aakritigkmit/payment-gateway/internal/handlers"
	"github.com/aakritigkmit/payment-gateway/internal/repository"
	"github.com/aakritigkmit/payment-gateway/internal/services"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupDBSRoutes(r chi.Router, db *mongo.Database) {
	dbsRepo := repository.NewDBSRepo(db)
	dbsService := services.NewDBSService(dbsRepo)
	dbsHandler := handlers.NewDBSHandler(dbsService)

	r.Post("/bank-statement", dbsHandler.HandleBankStatement)
}
