package dto

import "time"

type PinelabsOrderDetailResponse struct {
	OrderID                 string          `json:"order_id"`
	MerchantOrderReference  string          `json:"merchant_order_reference"`
	Type                    string          `json:"type"`
	Status                  string          `json:"status"`
	CallbackURL             string          `json:"callback_url"`
	FailureCallbackURL      string          `json:"failure_callback_url"`
	MerchantID              string          `json:"merchant_id"`
	OrderAmount             Amount          `json:"order_amount"`
	Notes                   string          `json:"notes"`
	PreAuth                 bool            `json:"pre_auth"`
	AllowedPaymentMethods   []string        `json:"allowed_payment_methods"`
	PurchaseDetails         PurchaseDetails `json:"purchase_details"`
	Payments                []Payment       `json:"payments"`
	IntegrationMode         string          `json:"integration_mode"`
	PaymentRetriesRemaining int             `json:"payment_retries_remaining"`
	CreatedAt               time.Time       `json:"created_at"`
	UpdatedAt               time.Time       `json:"updated_at"`
}

type Amount struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

type PurchaseDetails struct {
	Customer         Customer          `json:"customer"`
	MerchantMetadata map[string]string `json:"merchant_metadata"`
}

type Payment struct {
	ID                       string        `json:"id"`
	MerchantPaymentReference string        `json:"merchant_payment_reference"`
	Status                   string        `json:"status"`
	PaymentAmount            Amount        `json:"payment_amount"`
	PaymentMethod            string        `json:"payment_method"`
	PaymentOption            PaymentOption `json:"payment_option"`
	AcquirerData             AcquirerData  `json:"acquirer_data"`
	ErrorDetail              ErrorDetail   `json:"error_detail"`
	CreatedAt                time.Time     `json:"created_at"`
	UpdatedAt                time.Time     `json:"updated_at"`
}

type PaymentOption struct {
	NetbankingData *NetbankingData `json:"netbanking_data,omitempty"`
	UPIData        *UPIData        `json:"upi_data,omitempty"`
}

type NetbankingData struct {
	PayCode string `json:"pay_code"`
	TxnMode string `json:"txn_mode"`
}

type UPIData struct {
	TxnMode string `json:"txn_mode"`
}

type AcquirerData struct {
	ApprovalCode      string `json:"approval_code"`
	AcquirerReference string `json:"acquirer_reference"`
	RRN               string `json:"rrn"`
	IsAggregator      bool   `json:"is_aggregator"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
