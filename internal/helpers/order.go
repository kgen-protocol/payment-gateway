package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/aakritigkmit/payment-gateway/internal/dto"
)

// StructToMap converts a struct to a map[string]interface{}
// Returns error if marshalling/unmarshalling fails
func BuildOrderPayload(req dto.PlaceOrderRequest) ([]byte, error) {
	payload := map[string]interface{}{
		"merchant_order_reference": req.MerchantOrderReference,
		"order_amount": map[string]interface{}{
			"value":    req.OrderAmount.Value,
			"currency": req.OrderAmount.Currency,
		},
		"pre_auth":                req.PreAuth,
		"allowed_payment_methods": req.AllowedPaymentMethods,
		"notes":                   req.Notes,
		"callback_url":            req.CallbackURL,
		"failure_callback_url":    req.FailureCallbackURL,
		"purchase_details": map[string]interface{}{
			"customer": map[string]interface{}{
				"email_id":      req.PurchaseDetails.Customer.EmailID,
				"first_name":    req.PurchaseDetails.Customer.FirstName,
				"last_name":     req.PurchaseDetails.Customer.LastName,
				"customer_id":   req.PurchaseDetails.Customer.CustomerID,
				"mobile_number": req.PurchaseDetails.Customer.MobileNumber,
				"billing_address": map[string]interface{}{
					"address1": req.PurchaseDetails.Customer.BillingAddress.Address1,
					"address2": req.PurchaseDetails.Customer.BillingAddress.Address2,
					"address3": req.PurchaseDetails.Customer.BillingAddress.Address3,
					"pincode":  req.PurchaseDetails.Customer.BillingAddress.Pincode,
					"city":     req.PurchaseDetails.Customer.BillingAddress.City,
					"state":    req.PurchaseDetails.Customer.BillingAddress.State,
					"country":  req.PurchaseDetails.Customer.BillingAddress.Country,
				},
				"shipping_address": map[string]interface{}{
					"address1": req.PurchaseDetails.Customer.ShippingAddress.Address1,
					"address2": req.PurchaseDetails.Customer.ShippingAddress.Address2,
					"address3": req.PurchaseDetails.Customer.ShippingAddress.Address3,
					"pincode":  req.PurchaseDetails.Customer.ShippingAddress.Pincode,
					"city":     req.PurchaseDetails.Customer.ShippingAddress.City,
					"state":    req.PurchaseDetails.Customer.ShippingAddress.State,
					"country":  req.PurchaseDetails.Customer.ShippingAddress.Country,
				},
			},
			"merchant_metadata": req.PurchaseDetails.MerchantMetadata,
		},
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return jsonBytes, nil
}
