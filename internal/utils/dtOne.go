// package utils

// import (
// 	"encoding/base64"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/aakritigkmit/payment-gateway/internal/dto"
// )

// // const defaultLimit = 100

// func FetchProductsChunk(offset int, limit int) ([]dto.Product, error) {
// 	username := os.Getenv("DTONE_USERNAME")
// 	password := os.Getenv("DTONE_PASSWORD")
// 	baseURL := os.Getenv("DTONE_BASE_URL")

// 	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
// 	url := fmt.Sprintf("%s/v1/products?offset=%d&limit=%d", baseURL, offset, limit)

// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Set("Authorization", "Basic "+auth)
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{Timeout: 60 * time.Second}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to call DTOne API: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != 200 {
// 		body, _ := io.ReadAll(resp.Body)
// 		return nil, fmt.Errorf("DTOne API error: %s", string(body))
// 	}

// 	var products []dto.Product
// 	body, _ := io.ReadAll(resp.Body)
// 	if err := json.Unmarshal(body, &products); err != nil {
// 		return nil, fmt.Errorf("error decoding response: %w", err)
// 	}

//		return products, nil
//	}
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aakritigkmit/payment-gateway/internal/model"
)

func FetchProductsFromDTOneAPI(username, password string) ([]model.Product, error) {
	req, err := http.NewRequest("GET", "https://preprod-dvs-api.dtone.com/v1/products", nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Raw body:", string(body))
		return nil, fmt.Errorf("API error: %s", body)
	}

	var products []model.Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, err
	}
	return products, nil

}
