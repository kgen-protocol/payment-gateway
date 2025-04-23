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

// func (h *ProductHandler) CreateBulkProductTransaction(w http.ResponseWriter, r *http.Request) {
// 	var req dto.BulkTransactionRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	// Respond immediately to the user
// 	utils.SendSuccessResponse(w, http.StatusOK, "Your order is being processed in the background", nil)

// 	// Run background job independently
// 	go func() {
// 		bgCtx := context.Background()
// 		if err := h.service.CreateAndSaveBulkTransactions(bgCtx, req); err != nil {
// 			log.Printf("Background bulk transaction failed: %v", err)
// 		} else {
// 			log.Println("Background bulk transaction completed successfully")
// 		}
// 	}()
// }

func (h *ProductHandler) CreateBulkProductTransaction(w http.ResponseWriter, r *http.Request) {
	var req dto.BulkTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Generate a new orderId and save it to DB
	orderId, err := h.service.InitBulkProductTransaction(r.Context(), req)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create transaction entry")
		return
	}

	// Respond immediately
	utils.SendSuccessResponse(w, http.StatusOK, "Your order has been placed", map[string]string{"orderId": orderId})

	// Start async background processing
	go func() {
		bgCtx := context.Background()
		if err := h.service.ProcessBulkProductTransactionAsync(bgCtx, req, orderId); err != nil {
			log.Printf("Async processing failed for OrderID %s: %v", orderId, err)
		} else {
			log.Printf("Bulk transaction completed successfully for OrderID %s", orderId)
		}
	}()
}
