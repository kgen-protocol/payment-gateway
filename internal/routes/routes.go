package routes

import (
	"log"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

// routeRegistry holds the mapping of route initialization functions
var routeRegistry = map[string]func(r chi.Router, db *mongo.Database){
	"auth":     SetupAuthRoutes,
	"orders":   SetupOrderRoutes,
	"products": SetupProductRoutes,
	"dbs":      SetupDBSRoutes,
	"health":   SetupHealthRoutes,
}

// SetupRoutes initializes all application routes with /api prefix
func SetupRoutes(r *chi.Mux, db *mongo.Database) {
	apiRouter := chi.NewRouter()

	for routeName, setupFunc := range routeRegistry {
		log.Println("Registering route:", routeName) // Debugging log
		apiRouter.Route("/"+routeName, func(subRouter chi.Router) {
			setupFunc(subRouter, db)
		})
	}

	r.Mount("/api", apiRouter)
}
