package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderAmount struct {
	Value    float32 `json:"value"`
	Currency string  `json:"currency"`
}

type PurchaseDetails struct {
	Customer         Customer          `json:"customer"`
	MerchantMetadata map[string]string `json:"merchant_metadata"`
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
	OrderAmount            OrderAmount        `json:"order_amount"`
	PreAuth                bool               `json:"pre_auth"`
	AllowedPaymentMethods  []string           `json:"allowed_payment_methods"`
	Notes                  string             `json:"notes"`
	CallbackURL            string             `json:"callback_url"`
	FailureCallbackURL     string             `json:"failure_callback_url"`
	PurchaseDetails        PurchaseDetails    `json:"purchase_details"`
	PineOrderID            string             `json:"order_id" bson:"order_id"`
	Token                  string             `json:"token"`
	RedirectURL            string             `json:"redirect_url"`
	Status                 string             `json:"status"` // e.g., "PROCESSED"
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
