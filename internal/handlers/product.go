package handlers

import (
	"fmt"
	"net/http"

	"github.com/aakritigkmit/payment-gateway/internal/services"
	"github.com/aakritigkmit/payment-gateway/internal/utils"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service}
}

func (h *ProductHandler) SyncProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := h.service.SyncProducts(ctx)
	if err != nil {
		fmt.Println("err: ", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, "Products fetched and saved successfully", nil)
}
