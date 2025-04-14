package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	orderID := r.FormValue("order_id")
	if orderID == "" {
		http.Error(w, "Missing order_id", http.StatusBadRequest)
		return
	}

	fmt.Println("orderID: ", orderID)

	status := r.FormValue("status")
	if status == "" {
		http.Error(w, "Missing status", http.StatusBadRequest)
		return
	}

	// // Create payload with status from Plural
	// updatePayload := &dto.UpdateOrderPayload{
	// 	Status: status,
	// }

	// err := h.service.UpdateOrder(orderID, updatePayload)
	// if err != nil {
	// 	utils.SendErrorResponse(w, http.StatusInternalServerError, "r statusFailed to update orde")
	// 	return
	// }

	if err := h.service.SyncOrderDataFromPinelabs(r.Context(), orderID); err != nil {
		fmt.Println("error: ", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to sync order data")
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, "Order status updated successfully", nil)
}
