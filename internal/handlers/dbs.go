package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aakritigkmit/payment-gateway/internal/model"
	"github.com/aakritigkmit/payment-gateway/internal/services"
	"github.com/aakritigkmit/payment-gateway/internal/utils"
)

type DBSHandler struct {
	service *services.DBSService
}

func NewDBSHandler(service *services.DBSService) *DBSHandler {
	return &DBSHandler{
		service: service,
	}
}

func (h *DBSHandler) HandleBankStatement(w http.ResponseWriter, r *http.Request) {
	var req model.Camt053Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := h.service.ProcessBankStatement(req)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, "Bank statement fetched", resp)
}
