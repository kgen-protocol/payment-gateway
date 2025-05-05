package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/aakritigkmit/payment-gateway/internal/config"
	"github.com/aakritigkmit/payment-gateway/internal/dto"
	"github.com/aakritigkmit/payment-gateway/internal/model"
)

func FetchDTOneProducts(ctx context.Context, page, perPage int, filter dto.ProductSyncRequest) ([]model.Product, int, error) {
	cfg := config.GetConfig()

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.DtOneUsername, cfg.DtOnePassword)))

	params := url.Values{}
	params.Set("per_page", strconv.Itoa(perPage))
	params.Set("page", strconv.Itoa(page))

	if filter.ServiceID != 0 {
		params.Set("service_id", strconv.Itoa(filter.ServiceID))
	}
	if filter.CountryISOCode != "" {
		params.Set("country_iso_code", filter.CountryISOCode)
	}
	if filter.Type != "" {
		params.Set("type", filter.Type)
	}

	if filter.OperatorID != 0 {
		params.Set("operator_id", strconv.Itoa(filter.OperatorID))
	}

	fullURL := fmt.Sprintf("%s?%s", cfg.DtOneProductsURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, 0, fmt.Errorf("status: %d, body: %s", resp.StatusCode, body)
	}

	var products []model.Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, 0, fmt.Errorf("decode error: %w", err)
	}

	totalPagesStr := resp.Header.Get("X-Total-Pages")
	totalPages, err := strconv.Atoi(totalPagesStr)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid total pages header: %w", err)
	}

	return products, totalPages, nil
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
