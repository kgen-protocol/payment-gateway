package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aakritigkmit/payment-gateway/internal/config"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

type OrderAPIResponse struct {
	Token       string `json:"token"`
	OrderID     string `json:"order_id"`
	RedirectURL string `json:"redirect_url"`
}

func FetchAccessToken(ctx context.Context) (TokenResponse, error) {
	cfg := config.GetConfig()

	body := map[string]string{
		"client_id":     cfg.PinelabsClientID,
		"client_secret": cfg.PinelabsClientSecret,
		"grant_type":    cfg.PinelabsGrantType,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, cfg.PinelabsTokenURL, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return TokenResponse{}, errors.New("token fetch failed")
	}
	defer resp.Body.Close()

	var token TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return TokenResponse{}, err
	}

	return token, nil
}

func CreateOrderRequest(ctx context.Context, token string, jsonPayload []byte) (OrderAPIResponse, error) {
	cfg := config.GetConfig()

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.PinelabsOrderURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return OrderAPIResponse{}, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return OrderAPIResponse{}, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	var response OrderAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return OrderAPIResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check status *after* decoding
	if resp.StatusCode != http.StatusOK {
		return OrderAPIResponse{}, fmt.Errorf("order creation failed: status %d, message: %+v", resp.StatusCode, response)
	}

	return response, nil
}
