package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/aakritigkmit/payment-gateway/internal/dto"
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
	var req dto.CAMT053Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := h.service.ProcessBankStatement(req)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, "Bank statement processed", nil)
}

func (h *DBSHandler) HandleIntradayNotification(w http.ResponseWriter, r *http.Request) {
	var payload model.IntradayNotificationPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.ProcessIntradayNotification(payload); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, "Intraday Notification processed successfully", nil)
}

func (h *DBSHandler) HandleIncomingNotification(w http.ResponseWriter, r *http.Request) {
	var payload model.IncomingNotificationPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.ProcessIncomingNotification(payload); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, "Incoming Notification processed successfully", nil)
}
func (h *DBSHandler) HandleDBSEvent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	// Try Camt053Request (Bank Statement)
	var bank dto.CAMT053Request
	if err := json.Unmarshal(body, &bank); err == nil && bank.TxnEnqResponse.MessageType != "" {
		err := h.service.ProcessBankStatement(bank)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.SendSuccessResponse(w, http.StatusOK, "Bank statement processed successfully", nil)
		return
	}

	// Try IntradayNotificationPayload
	var intraday model.IntradayNotificationPayload
	if err := json.Unmarshal(body, &intraday); err == nil && intraday.TxnInfo.TxnType != "" {
		if err := h.service.ProcessIntradayNotification(intraday); err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.SendSuccessResponse(w, http.StatusOK, "Intraday Notification processed successfully", nil)
		return
	}

	// Try IncomingNotificationPayload
	var incoming model.IncomingNotificationPayload
	if err := json.Unmarshal(body, &incoming); err == nil && incoming.TxnInfo.TxnType != "" {
		if err := h.service.ProcessIncomingNotification(incoming); err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.SendSuccessResponse(w, http.StatusOK, "Incoming Notification processed successfully", nil)
		return
	}

	utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid or unrecognized payload")
}
