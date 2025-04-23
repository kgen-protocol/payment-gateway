package dto

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type UserRequest struct {
	FirstName      string `bson:"name,omitempty" validate:"required,min=2,max=100"`
	LastName       string `bson:"name,omitempty" validate:"required,min=2,max=100"`
	Email          string `bson:"email,omitempty" validate:"required,email"`
	Password       string `bson:"password,omitempty" json:"password" validate:"required,min=6"`
	MobileNumber   string `json:"mobile_number,omitempty"`
	BillingAddress string `json:"billing_address,omitempty"`
}

func (u *UserRequest) ValidateUser() error {
	validate := validator.New()
	return validate.Struct(u)
}

type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type PineTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type PineTokenResponse struct {
	Token string `json:"access_token"`
}

type PineCheckoutResponse struct {
	Token           string `json:"token"`
	OrderID         string `json:"order_id"`
	RedirectURL     string `json:"redirect_url"`
	ResponseCode    int    `json:"response_code"`
	ResponseMessage string `json:"response_message"`
}

type Order struct {
	UserID                 string  `json:"userId"`
	TransactionReferenceId string  `json:"transactionReferenceId"`
	Amount                 float32 `json:"amount"`
	Currency               string  `json:"currency"`
	Status                 string  `json:"status"` // pending, success, failure

}

type PlaceOrderRequest struct {
	MerchantOrderReference int64          `json:"merchant_order_reference"`
	OrderAmount            OrderAmount    `json:"order_amount"`
	PreAuth                bool           `json:"pre_auth"`
	AllowedPaymentMethods  []string       `json:"allowed_payment_methods"`
	Notes                  string         `json:"notes"`
	CallbackURL            string         `json:"callback_url"`
	FailureCallbackURL     string         `json:"failure_callback_url"`
	PurchaseDetails        PurchaseDetail `json:"purchase_details"`
}

type OrderAmount struct {
	Value    float32 `json:"value"`
	Currency string  `json:"currency"`
}

type Customer struct {
	EmailID         string  `json:"email_id"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	CustomerID      string  `json:"customer_id"`
	MobileNumber    string  `json:"mobile_number"`
	CountryCode     string  `json:"country_code"`
	Status          string  `json:"status"`
	BillingAddress  Address `json:"billing_address"`
	ShippingAddress Address `json:"shipping_address"`
	IsEditAllowed   bool    `json:"is_edit_customer_details_allowed"`
}

type Address struct {
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	Address3 string `json:"address3"`
	Pincode  string `json:"pincode"`
	City     string `json:"city"`
	State    string `json:"state"`
	Country  string `json:"country"`
}

type TransactionCallbackPayload struct {
	TransactionReferenceId string `json:"transactionReferenceId"`
}

type UpdateOrderPayload struct {
	Status   string  `json:"status,omitempty"`
	Amount   float32 `json:"amount,omitempty"`
	Currency string  `json:"currency,omitempty"`
	UserID   string  `json:"userId,omitempty"`
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

type AcquirerData struct {
	ApprovalCode      string `json:"approval_code" bson:"approval_code"`
	AcquirerReference string `json:"acquirer_reference" bson:"acquirer_reference"`
	RRN               string `json:"rrn" bson:"rrn"`
	IsAggregator      bool   `json:"is_aggregator" bson:"is_aggregator"`
	AcquirerName      string `json:"acquirer_name" bson:"acquirer_name"`
}

type PineOrderResponse struct {
	Data PineOrderData `json:"data"`
}

type PineOrderData struct {
	OrderID                 string         `json:"order_id"`
	MerchantOrderReference  string         `json:"merchant_order_reference"`
	Type                    string         `json:"type"`
	Status                  string         `json:"status"`
	CallbackURL             string         `json:"callback_url"`
	FailureCallbackURL      string         `json:"failure_callback_url"`
	MerchantID              string         `json:"merchant_id"`
	OrderAmount             OrderAmount    `json:"order_amount"`
	Notes                   string         `json:"notes"`
	PreAuth                 bool           `json:"pre_auth"`
	AllowedPaymentMethods   []string       `json:"allowed_payment_methods"`
	PurchaseDetails         PurchaseDetail `json:"purchase_details"`
	Payments                []Payment      `json:"payments"`
	Refunds                 []Refund       `json:"refunds"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
	IntegrationMode         string         `json:"integration_mode"`
	PaymentRetriesRemaining int            `json:"payment_retries_remaining"`
}

// type RefundRequest struct {
// 	OrderID string `json:"order_id"`
// }

// type RefundAPIResponse struct {
// 	Success bool                  `json:"success"`
// 	Message string                `json:"message"`
// 	Data    RefundAPIResponseData `json:"data"`
// }

// type RefundAPIResponseData struct {
// 	Status   string `json:"status"`
// 	RefundID string `json:"refund_id"`
// }

// type RefundPayload struct {
// 	OrderID                string           `json:"order_id"`
// 	MerchantOrderReference string           `json:"merchant_order_reference"`
// 	OrderAmount            OrderAmount      `json:"order_amount"`
// 	Amount                 interface{}      `json:"amount"`
// 	MerchantMetadata       MerchantMetadata `json:"merchant_meta_data"`
// }
// type MerchantMetadata struct {
// 	Key1 string `json:"key1"`
// 	Key2 string `json:"key_2"`
// }

type RefundRequest struct {
	OrderID string `json:"order_id"`
	// MerchantOrderReference string           `json:"merchant_order_reference"`
	OrderAmount int `json:"order_amount"`
	// MerchantMetadata       MerchantMetadata `json:"merchant_metadata"`
}

type MerchantMetadata struct {
	Key1 string `json:"key1"`
	Key2 string `json:"key_2"`
}

type RefundResponse struct {
	Data struct {
		OrderID                string           `json:"order_id"`
		MerchantOrderReference string           `json:"merchant_order_reference"`
		Type                   string           `json:"type"`
		Status                 string           `json:"status"`
		CallbackURL            string           `json:"callback_url"`
		FailureCallbackURL     string           `json:"failure_callback_url"`
		MerchantID             string           `json:"merchant_id"`
		OrderAmount            OrderAmount      `json:"order_amount"`
		Notes                  string           `json:"notes"`
		PreAuth                bool             `json:"pre_auth"`
		AllowedPaymentMethods  []string         `json:"allowed_payment_methods"`
		PurchaseDetails        PurchaseDetail   `json:"purchase_details"`
		Payments               []Payment        `json:"payments"`
		Refunds                []Refund         `json:"refunds"`
		MerchantMetadata       MerchantMetadata `json:"merchant_meta_data"`
		CreatedAt              string           `json:"created_at"`
		UpdatedAt              string           `json:"updated_at"`
		IntegrationMode        string           `json:"integration_mode"`
	} `json:"data"`
}

type PurchaseDetail struct {
	Customer         Customer         `json:"customer"`
	MerchantMetadata MerchantMetadata `json:"merchant_metadata"`
}

type Refund struct {
	MerchantOrderReference string         `json:"merchant_order_reference"`
	OrderID                string         `json:"order_id"`
	Type                   string         `json:"type"`
	Status                 string         `json:"status"`
	OrderAmount            OrderAmount    `json:"order_amount"`
	Payments               []Payment      `json:"payments"`
	PurchaseDetails        PurchaseDetail `json:"purchase_details"`
	CreatedAt              string         `json:"created_at"`
	UpdatedAt              string         `json:"updated_at"`
}

type RefundOrderResponse struct {
	Data RefundOrderData `json:"data"`
}

type RefundOrderData struct {
	OrderID                 string          `json:"order_id"`
	ParentOrderID           string          `json:"parent_order_id"`
	MerchantOrderReference  string          `json:"merchant_order_reference"`
	Type                    string          `json:"type"`
	Status                  string          `json:"status"`
	MerchantID              string          `json:"merchant_id"`
	OrderAmount             Amount          `json:"order_amount"`
	PurchaseDetails         PurchaseDetails `json:"purchase_details"`
	Payments                []Payment       `json:"payments"`
	CreatedAt               string          `json:"created_at"`
	UpdatedAt               string          `json:"updated_at"`
	IntegrationMode         string          `json:"integration_mode"`
	PaymentRetriesRemaining int             `json:"payment_retries_remaining"`
}

type PurchaseDetails struct {
	Customer         Customer          `json:"customer"`
	MerchantMetadata map[string]string `json:"merchant_metadata"`
}
