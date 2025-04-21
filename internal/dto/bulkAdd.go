package dto

type LineItem struct {
	ProductID int     `json:"productId"`
	Quantity  int     `json:"quantity"`
	Amount    float64 `json:"amount"`
}

type BulkTransactionRequest struct {
	LineItems    []LineItem `json:"lineItems"`
	MobileNumber string     `json:"mobile_number"`
}
