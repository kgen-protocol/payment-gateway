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
	var payload dto.TransactionCallbackPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatePayload := &dto.UpdateOrderPayload{
		Status: "success",
	}

	err := h.service.UpdateOrder(payload.TransactionReferenceId, updatePayload)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update order status")
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, "Order status updated successfully", nil)
}

func (h *OrderHandler) HandleFailureCallback(w http.ResponseWriter, r *http.Request) {
	var payload dto.TransactionCallbackPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	updatePayload := &dto.UpdateOrderPayload{
		Status: "failed",
	}
	err := h.service.UpdateOrder(payload.TransactionReferenceId, updatePayload)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update order status")
		return
	}
	utils.SendSuccessResponse(w, http.StatusOK, "Order status updated successfully", nil)
}
