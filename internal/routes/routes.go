package routes

import (
	"log"
	"net/http"

	"github.com/aakritigkmit/payment-gateway/internal/utils"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

// routeRegistry holds the mapping of route initialization functions
var routeRegistry = map[string]func(r chi.Router, db *mongo.Database){
	"auth":             SetupAuthRoutes,
	"orders":           SetupOrderRoutes,
	"products":         SetupProductRoutes,
	"payments-webhook": SetupDBSRoutes,
}

// SetupRoutes initializes all application routes with /api prefix
func SetupRoutes(r *chi.Mux, db *mongo.Database) {

	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		utils.SendSuccessResponse(w, http.StatusOK, "API is working correctly", nil)
	})

	apiRouter := chi.NewRouter()

	for routeName, setupFunc := range routeRegistry {
		log.Println("Registering route:", routeName) // Debugging log
		apiRouter.Route("/"+routeName, func(subRouter chi.Router) {
			setupFunc(subRouter, db)
		})
	}

	r.Mount("/api/v1", apiRouter)
}
