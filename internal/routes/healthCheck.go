package routes

import (
	"net/http"

	"github.com/aakritigkmit/payment-gateway/internal/utils"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupHealthRoutes(r chi.Router, _ *mongo.Database) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.SendSuccessResponse(w, http.StatusOK, "API is working correctly", nil)
	})
}
