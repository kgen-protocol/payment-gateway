package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aakritigkmit/payment-gateway/internal/dto"
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

func (h *ProductHandler) HandleProductTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.CreateAndSaveTransaction(ctx, req); err != nil {
		log.Printf("Failed: %v", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, "Transaction created and saved successfully", nil)
}

func (h *ProductHandler) CreateBulkProductTransaction(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	var req dto.BulkTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	go func(request dto.BulkTransactionRequest) {
		// Create a detached context so work isn't canceled if client disconnects
		bgCtx := context.Background()
		if err := h.service.CreateAndSaveBulkTransactions(bgCtx, request); err != nil {
			log.Printf("Async bulk transaction failed: %v", err)
		}
	}(req)

	utils.SendSuccessResponse(w, http.StatusOK, "Bulk transactions processed successfully", nil)
}
