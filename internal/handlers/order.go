package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aakritigkmit/payment-gateway/internal/config"
	"github.com/aakritigkmit/payment-gateway/internal/dto"
	"github.com/aakritigkmit/payment-gateway/internal/services"
	"github.com/aakritigkmit/payment-gateway/internal/utils"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{service}
}

func (h *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.PlaceOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	orderResp, err := h.service.PlaceOrder(r.Context(), req)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, "Order placed successfully", map[string]string{
		"token":        orderResp.Token,
		"order_id":     orderResp.OrderID,
		"redirect_url": orderResp.RedirectURL,
	})
}

func (h *OrderHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {

	cfg := config.GetConfig()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	orderID := r.FormValue("order_id")
	if orderID == "" {
		http.Error(w, "Missing order_id", http.StatusBadRequest)
		return
	}

	status := r.FormValue("status")
	if status == "" {
		http.Error(w, "Missing status", http.StatusBadRequest)
		return
	}
	fmt.Println("status", status)

	receivedSignature := r.FormValue("signature")
	if receivedSignature == "" {
		http.Error(w, "Missing signature", http.StatusBadRequest)
		return
	}
	fmt.Println("signature", receivedSignature)

	errorCode := r.FormValue("error_code")
	errorMessage := r.FormValue("error_message")

	// --- MANDATORY SIGNATURE VERIFICATION STEP ---
	// Generate the signature on your server using the same parameters and secret key.
	generatedSignature, err := utils.GenerateServerSignature(orderID, status, errorCode, errorMessage, cfg.PinelabsClientSecret)
	fmt.Printf("error-generated signature: %s\n", generatedSignature)
	if err != nil {
		fmt.Printf("Error generating server signature: %v\n", err)
		http.Error(w, "Internal server error during signature generation", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Server-generated signature: '%s'\n", generatedSignature)

	// Compare the generated signature with the received signature.
	if strings.ToUpper(generatedSignature) != strings.ToUpper(receivedSignature) {
		fmt.Printf("Signature mismatch! Received: '%s', Generated: '%s'\n", receivedSignature, generatedSignature)
		http.Error(w, "Invalid signature: Callback authenticity could not be verified", http.StatusUnauthorized)
		return
	}
	fmt.Println("Signature verified successfully! Proceeding with transaction update.")

	go h.service.FetchAndUpdateTransactionDetails(context.Background(), orderID, receivedSignature)

	utils.SendSuccessResponse(w, http.StatusOK, "Order status updated successfully", nil)
}

func (h *OrderHandler) RefundOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.RefundRequest
	fmt.Println("req", req)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	refundResp, err := h.service.ProcessRefund(r.Context(), req)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, "Refund processed successfully", refundResp)
}
