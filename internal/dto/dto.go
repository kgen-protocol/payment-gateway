package dto

import (
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

type PurchaseDetail struct {
	Customer         Customer          `json:"customer"`
	MerchantMetadata map[string]string `json:"merchant_metadata"`
}

type Customer struct {
	EmailID         string  `json:"email_id"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	CustomerID      string  `json:"customer_id"`
	MobileNumber    string  `json:"mobile_number"`
	BillingAddress  Address `json:"billing_address"`
	ShippingAddress Address `json:"shipping_address"`
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
