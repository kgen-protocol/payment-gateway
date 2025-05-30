package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aakritigkmit/payment-gateway/internal/config"
	"github.com/aakritigkmit/payment-gateway/internal/dto"
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
	const redisTokenKey = "pinelabs:access_token"

	// Try fetching from Redis first
	cachedToken, err := GetRedisKey(ctx, redisTokenKey)
	if err == nil && cachedToken != "" {
		log.Println("Using cached Pinelabs access token")
		return TokenResponse{AccessToken: cachedToken}, nil
	}

	cfg := config.GetConfig()

	// Make API call if not in Redis
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

	// Save token in Redis (set to expire in 55 minutes)
	if err := SetRedisKey(ctx, redisTokenKey, token.AccessToken, 55*time.Minute); err != nil {
		log.Printf("Failed to cache access token in Redis: %v", err)
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

func GetOrderDetails(ctx context.Context, token string, pineOrderID string) (*dto.PineOrderResponse, error) {
	cfg := config.GetConfig()

	// Build request URL
	url := fmt.Sprintf("%s/%s", cfg.PinelabsGetOrderURL, pineOrderID)

	// Build request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch order details: status %d", resp.StatusCode)
	}

	// Decode response into your DTO
	var result dto.PineOrderResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func CreateRefundRequest(ctx context.Context, accessToken, orderID string, payload []byte) (*dto.RefundOrderResponse, error) {
	cfg := config.GetConfig()
	url := fmt.Sprintf("%s/%s", cfg.PinelabsRefundURL, orderID)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("refund API error: %s", string(body))
	}

	var refundResp dto.RefundOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&refundResp); err != nil {
		return nil, err
	}

	return &refundResp, nil
}
