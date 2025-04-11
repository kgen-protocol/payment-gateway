package handlers

import (
	"encoding/json"
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

func (h *OrderHandler) HandleSuccessCallback(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	orderID := r.FormValue("order_id")
	if orderID == "" {
		http.Error(w, "Missing order_id in callback", http.StatusBadRequest)
		return
	}
	updatePayload := &dto.UpdateOrderPayload{
		Status: "success",
	}

	err := h.service.UpdateOrder(orderID, updatePayload)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update order status")
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, "Order status updated successfully", nil)
}

func (h *OrderHandler) HandleFailureCallback(w http.ResponseWriter, r *http.Request) {
	// 1. Try reading from POST body (assuming Pine sends JSON)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	orderID := r.FormValue("order_id")
	if orderID == "" {
		http.Error(w, "Missing order_id in callback", http.StatusBadRequest)
		return
	}
	// 3. Update order status to "failed"
	updatePayload := &dto.UpdateOrderPayload{
		Status: "failed",
	}

	if err := h.service.UpdateOrder(orderID, updatePayload); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update order status")
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, "Order status updated successfully", nil)
}
