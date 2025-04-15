package handlers

import (
	"fmt"
	"net/http"

	"github.com/aakritigkmit/payment-gateway/internal/services"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service}
}

func (h *ProductHandler) SyncProductsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := h.service.SyncProducts(ctx)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to sync products: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("✅ Product sync completed"))
}
