package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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

func CreateDTOneTransaction(ctx context.Context, externalID string, productID int, mobileNumber string) error {
	cfg := config.GetConfig()

	payload := map[string]interface{}{
		"external_id": externalID,
		"product_id":  productID,
		"credit_party_identifier": map[string]string{
			"mobile_number": mobileNumber,
		},
		"auto_confirm": true,
	}
	data, _ := json.Marshal(payload)

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.DtOneUsername, cfg.DtOnePassword)))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.DtOneTransactionURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create transaction, status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func FetchDTOneTransactionByExternalID(ctx context.Context, externalID string) ([]model.ProductTransaction, error) {
	cfg := config.GetConfig()
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.DtOneUsername, cfg.DtOnePassword)))
	url := fmt.Sprintf("%s?external_id=%s", cfg.DtOneGetTransactionURL, externalID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var txs []model.ProductTransaction
	if err := json.NewDecoder(resp.Body).Decode(&txs); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch transaction: status %d", resp.StatusCode)
	}

	return txs, nil
}
