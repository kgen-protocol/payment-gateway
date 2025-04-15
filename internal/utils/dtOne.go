package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aakritigkmit/payment-gateway/internal/config"
	"github.com/aakritigkmit/payment-gateway/internal/model"
)

func FetchDTOneProducts(ctx context.Context) ([]model.Product, error) {
	cfg := config.GetConfig()

	// Construct Basic Auth
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.DtOneUsername, cfg.DtOnePassword)))

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.DtOneProductsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+auth)

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Decode response body
	var products []model.Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch products: status %d, response: %+v", resp.StatusCode, products)
	}

	return products, nil
}
