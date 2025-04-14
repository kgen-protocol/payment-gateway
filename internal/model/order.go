package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type Order struct {
// 	ID                     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	UserID                 string             `bson:"userId"`
// 	TransactionReferenceId string             `bson:"transactionReferenceId"`
// 	// PineToken     string             `bson:"pineToken"`
// 	Amount    float32   `bson:"amount"`
// 	Currency  string    `bson:"currency"`
// 	Status    string    `bson:"status"` // pending, success, failure
// 	CreatedAt time.Time `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
// 	UpdatedAt time.Time `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
// }

type Order struct {
	ID                      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	OrderID                 string             `bson:"order_id"`
	MerchantOrderReference  string             `bson:"merchant_order_reference"`
	Type                    string             `bson:"type"`
	Status                  string             `bson:"status"`
	CallbackURL             string             `bson:"callback_url"`
	FailureCallbackURL      string             `bson:"failure_callback_url"`
	MerchantID              string             `bson:"merchant_id"`
	Amount                  OrderAmount        `bson:"amount"`
	Notes                   string             `bson:"notes"`
	PreAuth                 bool               `bson:"pre_auth"`
	AllowedPaymentMethods   []string           `bson:"allowed_payment_methods"`
	PurchaseDetails         PurchaseDetails    `bson:"purchase_details"`
	IntegrationMode         string             `bson:"integration_mode"`
	PaymentRetriesRemaining int                `bson:"payment_retries_remaining"`
	CreatedAt               time.Time          `bson:"created_at,omitempty"`
	UpdatedAt               time.Time          `bson:"updated_at,omitempty"`
}
