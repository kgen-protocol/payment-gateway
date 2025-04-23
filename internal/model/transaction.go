package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderAmount struct {
	Value    float32 `json:"value"`
	Currency string  `json:"currency"`
}

type PurchaseDetail struct {
	Customer         Customer         `json:"customer"`
	MerchantMetadata MerchantMetadata `json:"merchant_metadata"`
}

type Customer struct {
	EmailID                      string  `json:"email_id" bson:"email_id"`
	FirstName                    string  `json:"first_name" bson:"first_name"`
	LastName                     string  `json:"last_name" bson:"last_name"`
	CustomerID                   string  `json:"customer_id" bson:"customer_id"`
	MobileNumber                 string  `json:"mobile_number" bson:"mobile_number"`
	CountryCode                  string  `json:"country_code" bson:"country_code"`
	BillingAddress               Address `json:"billing_address" bson:"billing_address"`
	ShippingAddress              Address `json:"shipping_address" bson:"shipping_address"`
	IsEditCustomerDetailsAllowed bool    `json:"is_edit_customer_details_allowed" bson:"is_edit_customer_details_allowed"`
}

type Transaction struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	OrderId                string             `bson:"orderId"`
	MerchantOrderReference string             `json:"merchant_order_reference"`
	Type                   string             `json:"type"`
	MerchantID             string             `json:"merchant_id"`
	OrderAmount            OrderAmount        `json:"order_amount"`
	PreAuth                bool               `json:"pre_auth"`
	AllowedPaymentMethods  []string           `json:"allowed_payment_methods"`
	Notes                  string             `json:"notes"`
	CallbackURL            string             `json:"callback_url"`
	FailureCallbackURL     string             `json:"failure_callback_url"`
	PurchaseDetails        PurchaseDetail     `json:"purchase_details"`
	PineOrderID            string             `json:"order_id" bson:"order_id"`
	Token                  string             `json:"token"`
	RedirectURL            string             `json:"redirect_url"`
	Status                 string             `json:"status"` // e.g., "PROCESSED"
	Refunds                []Refund           `json:"refunds" bson:"refunds"`
	IntegrationMode        string             `json:"integration_mode" bson:"integration_mode"`
	Payments               []Payment          `json:"payments" bson:"payments"`
	CreatedAt              time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt              time.Time          `json:"updated_at" bson:"updated_at"`
}

type AcquirerData struct {
	ApprovalCode      string `json:"approval_code" bson:"approval_code"`
	AcquirerReference string `json:"acquirer_reference" bson:"acquirer_reference"`
	RRN               string `json:"rrn" bson:"rrn"`
	IsAggregator      bool   `json:"is_aggregator" bson:"is_aggregator"`
	AcquirerName      string `json:"acquirer_name" bson:"acquirer_name"`
}

type Payment struct {
	ID                       string       `json:"id" bson:"id"`
	MerchantPaymentReference string       `json:"merchant_payment_reference" bson:"merchant_payment_reference"`
	Status                   string       `json:"status" bson:"status"`
	PaymentAmount            OrderAmount  `json:"payment_amount" bson:"payment_amount"`
	PaymentMethod            string       `json:"payment_method" bson:"payment_method"`
	AcquirerData             AcquirerData `json:"acquirer_data" bson:"acquirer_data"`
	CreatedAt                time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt                time.Time    `json:"updated_at" bson:"updated_at"`
}

type Refund struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	TransactionID          primitive.ObjectID `bson:"transaction_id" json:"transaction_id"`
	MerchantOrderReference string             `json:"merchant_order_reference"`
	OrderID                string             `json:"order_id"`
	Type                   string             `json:"type"`
	Status                 string             `json:"status"`
	OrderAmount            OrderAmount        `json:"order_amount"`
	Payments               []Payment          `json:"payments"`
	PurchaseDetails        PurchaseDetail     `json:"purchase_details"`
	CreatedAt              string             `json:"created_at"`
	UpdatedAt              string             `json:"updated_at"`
}

type MerchantMetadata struct {
	Key1 string `json:"key1"`
	Key2 string `json:"key_2"`
}

type RefundOrderResponse struct {
	Data RefundOrderData `bson:"data" json:"data"`
}

type RefundOrderData struct {
	OrderID                 string          `bson:"order_id" json:"order_id"`
	ParentOrderID           string          `bson:"parent_order_id" json:"parent_order_id"`
	MerchantOrderReference  string          `bson:"merchant_order_reference" json:"merchant_order_reference"`
	Type                    string          `bson:"type" json:"type"`
	Status                  string          `bson:"status" json:"status"`
	MerchantID              string          `bson:"merchant_id" json:"merchant_id"`
	OrderAmount             Amount          `bson:"order_amount" json:"order_amount"`
	PurchaseDetails         PurchaseDetails `bson:"purchase_details" json:"purchase_details"`
	Payments                []Payment       `bson:"payments" json:"payments"`
	CreatedAt               string          `bson:"created_at" json:"created_at"`
	UpdatedAt               string          `bson:"updated_at" json:"updated_at"`
	IntegrationMode         string          `bson:"integration_mode" json:"integration_mode"`
	PaymentRetriesRemaining int             `bson:"payment_retries_remaining" json:"payment_retries_remaining"`
}

type PurchaseDetails struct {
	Customer         Customer          `bson:"customer" json:"customer"`
	MerchantMetadata map[string]string `bson:"merchant_metadata" json:"merchant_metadata"`
}
