package handlers

import (
	"context"
	"encoding/json"
	"io"
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

	var req dto.ProductSyncRequest
	if r.Body != nil {
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil && err != io.EOF {
			utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
			return
		}
	}

	// Launch background job
	go func() {
		bgCtx := context.Background()

		err := h.service.SyncProducts(bgCtx)
		if err != nil {
			log.Printf("Background sync failed: %v", err)
		} else {
			log.Println("Background sync completed successfully.")
		}
	}()

	utils.SendSuccessResponse(w, http.StatusAccepted, "Products are syncing in the background", nil)
}

func (h *ProductHandler) GenerateProductReport(w http.ResponseWriter, r *http.Request) {
	var req dto.ProductSyncRequest
	if r.Body != nil {
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
			utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
			return
		}
	}

	go func() {
		bgCtx := context.Background()
		if err := h.service.GenerateProductFetchReport(bgCtx, req); err != nil {
			log.Printf("Product report generation failed: %v", err)
		} else {
			log.Println("Product report generated successfully.")
		}
	}()

	utils.SendSuccessResponse(w, http.StatusAccepted, "Product report is being generated in the background", nil)
}

func (h *ProductHandler) GenerateProductReportByIDs(w http.ResponseWriter, r *http.Request) {
	var req dto.ProductReportByIDRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.ProductIDs) == 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid or missing product IDs")
		return
	}

	go func() {
		bgCtx := context.Background()
		if err := h.service.GenerateProductReportByIDs(bgCtx, req.ProductIDs); err != nil {
			log.Printf("[ReportByIDs] Report generation failed: %v", err)
		} else {
			log.Println("[ReportByIDs] Report generated successfully.")
		}
	}()

	utils.SendSuccessResponse(w, http.StatusAccepted, "Product report is being generated in the background", nil)
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
