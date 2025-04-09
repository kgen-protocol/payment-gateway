package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aakritigkmit/payment-gateway/internal/dto"
	"github.com/aakritigkmit/payment-gateway/internal/services"
	"github.com/aakritigkmit/payment-gateway/internal/utils"
)

type AuthHandler struct {
	Service *services.AuthService
}

func NewAuthHandler(services *services.AuthService) *AuthHandler {
	return &AuthHandler{
		Service: services,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var userReq dto.UserRequest

	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := userReq.ValidateUser(); err != nil {
		fmt.Println("hello")
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.Service.RegisterUser(userReq)
	if err != nil {
		fmt.Println("Error creating user:", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error creating user")
		return
	}
	utils.SendSuccessResponse(w, http.StatusCreated, "User registered successfully", nil)

}

// // Login an existing user
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials dto.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	token, err := h.Service.Login(credentials.Email, credentials.Password)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	utils.SendSuccessResponse(w, http.StatusOK, "Login successful", map[string]string{"token": token})
}
