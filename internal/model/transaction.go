package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderAmount struct {
	Value    float32 `json:"value"`
	Currency string  `json:"currency"`
}

type PurchaseDetails struct {
	Customer         Customer          `json:"customer"`
	MerchantMetadata map[string]string `json:"merchant_metadata"`
}

type Customer struct {
	EmailID         string  `json:"email_id" bson:"email_id"`
	FirstName       string  `json:"first_name" bson:"first_name"`
	LastName        string  `json:"last_name" bson:"last_name"`
	CustomerID      string  `json:"customer_id" bson:"customer_id"`
	MobileNumber    string  `json:"mobile_number" bson:"mobile_number"`
	BillingAddress  Address `json:"billing_address" bson:"billing_address"`
	ShippingAddress Address `json:"shipping_address" bson:"shipping_address"`
}
type Transaction struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	OrderId                string             `bson:"orderId"`
	MerchantOrderReference int64              `json:"merchant_order_reference"`
	OrderAmount            OrderAmount        `json:"order_amount"`
	PreAuth                bool               `json:"pre_auth"`
	AllowedPaymentMethods  []string           `json:"allowed_payment_methods"`
	Notes                  string             `json:"notes"`
	CallbackURL            string             `json:"callback_url"`
	FailureCallbackURL     string             `json:"failure_callback_url"`
	PurchaseDetails        PurchaseDetails    `json:"purchase_details"`
}
