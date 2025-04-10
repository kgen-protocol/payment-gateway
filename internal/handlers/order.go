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
		"order_id":     orderResp.OrderID,
		"redirect_url": orderResp.RedirectURL,
	})
}

func (h *OrderHandler) HandleSuccessCallback(w http.ResponseWriter, r *http.Request) {
	// You may want to read order ID or merchant reference from query/form data
	h.service.HandleSuccess(r.Context())
	utils.SendSuccessResponse(w, http.StatusOK, "Success callback received", nil)
}

func (h *OrderHandler) HandleFailureCallback(w http.ResponseWriter, r *http.Request) {
	h.service.HandleFailure(r.Context())
	utils.SendSuccessResponse(w, http.StatusOK, "Failure callback received", nil)
}
