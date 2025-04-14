package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderAmount struct {
	Value    float64 `json:"value" bson:"value"`
	Currency string  `json:"currency" bson:"currency"`
}

type PurchaseDetails struct {
	Customer         Customer          `json:"customer" bson:"customer"`
	MerchantMetadata map[string]string `json:"merchant_metadata" bson:"merchant_metadata"`
}

type Customer struct {
	EmailID                      string  `json:"email_id" bson:"email_id"`
	FirstName                    string  `json:"first_name" bson:"first_name"`
	LastName                     string  `json:"last_name" bson:"last_name"`
	CustomerID                   string  `json:"customer_id" bson:"customer_id"`
	MobileNumber                 string  `json:"mobile_number" bson:"mobile_number"`
	CountryCode                  string  `json:"country_code" bson:"country_code"`
	IsEditCustomerDetailsAllowed bool    `json:"is_edit_customer_details_allowed" bson:"is_edit_customer_details_allowed"`
	BillingAddress               Address `json:"billing_address" bson:"billing_address"`
	ShippingAddress              Address `json:"shipping_address" bson:"shipping_address"`
}

// type Transaction struct {
// 	ID                     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	OrderId                string             `bson:"orderId"`
// 	MerchantOrderReference int64              `json:"merchant_order_reference"`
// 	OrderAmount            OrderAmount        `json:"order_amount"`
// 	PreAuth                bool               `json:"pre_auth"`
// 	AllowedPaymentMethods  []string           `json:"allowed_payment_methods"`
// 	Notes                  string             `json:"notes"`
// 	CallbackURL            string             `json:"callback_url"`
// 	FailureCallbackURL     string             `json:"failure_callback_url"`
// 	PurchaseDetails        PurchaseDetails    `json:"purchase_details"`
// 	PineOrderID            string             `json:"pine_order_id"`
// 	Token                  string             `json:"token"`
// 	RedirectURL            string             `json:"redirect_url"`
// }

type Transaction struct {
	ID                       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	TransactionID            string             `bson:"transaction_id"`
	OrderID                  string             `bson:"order_id"`
	MerchantPaymentReference string             `bson:"merchant_payment_reference"`
	Status                   string             `bson:"status"`
	PaymentAmount            OrderAmount        `bson:"payment_amount"`
	PaymentMethod            string             `bson:"payment_method"`
	PaymentOption            PaymentOption      `bson:"payment_option"`
	AcquirerData             AcquirerData       `bson:"acquirer_data"`
	ErrorDetail              ErrorDetail        `bson:"error_detail"`
	CreatedAt                time.Time          `bson:"created_at,omitempty"`
	UpdatedAt                time.Time          `bson:"updated_at,omitempty"`
}

type PaymentOption struct {
	NetbankingData *NetbankingData `json:"netbanking_data,omitempty" bson:"netbanking_data,omitempty"`
	UPIData        *UPIData        `json:"upi_data,omitempty" bson:"upi_data,omitempty"`
}

type NetbankingData struct {
	PayCode string `json:"pay_code" bson:"pay_code"`
	TxnMode string `json:"txn_mode" bson:"txn_mode"`
}

type UPIData struct {
	TxnMode string `json:"txn_mode" bson:"txn_mode"`
}

type AcquirerData struct {
	ApprovalCode      string `json:"approval_code" bson:"approval_code"`
	AcquirerReference string `json:"acquirer_reference" bson:"acquirer_reference"`
	RRN               string `json:"rrn" bson:"rrn"`
	IsAggregator      bool   `json:"is_aggregator" bson:"is_aggregator"`
}

type ErrorDetail struct {
	Code    string `json:"code" bson:"code"`
	Message string `json:"message" bson:"message"`
}
