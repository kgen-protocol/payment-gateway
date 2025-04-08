package utils

import (
	"encoding/json"
	"net/http"
)

// JSONResponse represents a standard API response format
type JSONResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SendJSONResponse sends a JSON response with a success flag, data, and a message
func SendJSONResponse(w http.ResponseWriter, status int, success bool, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := JSONResponse{
		Success: success,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

// SendSuccessResponse sends a success JSON response
func SendSuccessResponse(w http.ResponseWriter, data interface{}, message string) {
	SendJSONResponse(w, http.StatusOK, true, data, message)
}

// SendErrorResponse sends an error JSON response
func SendErrorResponse(w http.ResponseWriter, status int, message string) {
	SendJSONResponse(w, status, false, nil, message)
}
