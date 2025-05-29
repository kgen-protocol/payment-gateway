package routes

import (
	"github.com/aakritigkmit/payment-gateway/internal/handlers"
	middlewares "github.com/aakritigkmit/payment-gateway/internal/middleware"
	"github.com/aakritigkmit/payment-gateway/internal/repository"
	"github.com/aakritigkmit/payment-gateway/internal/services"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupDBSRoutes(r chi.Router, db *mongo.Database) {
	dbsRepo := repository.NewDBSRepo(db)
	dbsService := services.NewDBSService(dbsRepo)
	dbsHandler := handlers.NewDBSHandler(dbsService)

	r.With(middlewares.AuthMiddleware).Post("/bank-statement", dbsHandler.HandleBankStatement)
	r.With(middlewares.AuthMiddleware).Post("/intraday/notification", dbsHandler.HandleIntradayNotification)
	r.With(middlewares.AuthMiddleware).Post("/incoming/notification", dbsHandler.HandleIncomingNotification)
	r.With(middlewares.AuthMiddleware).Post("/dbs", dbsHandler.HandleDBSEvent)

}
