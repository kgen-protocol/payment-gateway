package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/aakritigkmit/payment-gateway/internal/dto"
	"github.com/aakritigkmit/payment-gateway/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func MapPineOrderToTransactionModel(data *dto.PineOrderResponse, signature string) model.Transaction {
	payments := make([]model.Payment, 0)

	for _, p := range data.Data.Payments {
		payments = append(payments, model.Payment{
			ID:                       p.ID,
			MerchantPaymentReference: p.MerchantPaymentReference,
			Status:                   p.Status,
			PaymentMethod:            p.PaymentMethod,
			PaymentAmount: model.OrderAmount{
				Value:    p.PaymentAmount.Value,
				Currency: p.PaymentAmount.Currency,
			},
			AcquirerData: model.AcquirerData{
				ApprovalCode:      p.AcquirerData.ApprovalCode,
				AcquirerReference: p.AcquirerData.AcquirerReference,
				RRN:               p.AcquirerData.RRN,
				IsAggregator:      p.AcquirerData.IsAggregator,
				AcquirerName:      p.AcquirerData.AcquirerName,
			},
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})

	}

	return model.Transaction{
		OrderId:                data.Data.OrderID,
		MerchantOrderReference: data.Data.MerchantOrderReference,
		OrderAmount: model.OrderAmount{
			Value:    data.Data.OrderAmount.Value,
			Currency: data.Data.OrderAmount.Currency,
		},
		PreAuth:               data.Data.PreAuth,
		AllowedPaymentMethods: data.Data.AllowedPaymentMethods,
		Notes:                 data.Data.Notes,
		Type:                  data.Data.Type,
		MerchantID:            data.Data.MerchantID,
		CallbackURL:           data.Data.CallbackURL,
		FailureCallbackURL:    data.Data.FailureCallbackURL,
		PurchaseDetails: model.PurchaseDetail{
			Customer: model.Customer{
				EmailID:         data.Data.PurchaseDetails.Customer.EmailID,
				FirstName:       data.Data.PurchaseDetails.Customer.FirstName,
				LastName:        data.Data.PurchaseDetails.Customer.LastName,
				CustomerID:      data.Data.PurchaseDetails.Customer.CustomerID,
				MobileNumber:    data.Data.PurchaseDetails.Customer.MobileNumber,
				BillingAddress:  model.Address(data.Data.PurchaseDetails.Customer.BillingAddress),
				ShippingAddress: model.Address(data.Data.PurchaseDetails.Customer.ShippingAddress),
			},
			MerchantMetadata: model.MerchantMetadata(data.Data.PurchaseDetails.MerchantMetadata),
		},
		PineOrderID:     data.Data.OrderID,
		Status:          data.Data.Status,
		Signature:       signature,
		IntegrationMode: data.Data.IntegrationMode,
		Payments:        payments,
		CreatedAt:       data.Data.CreatedAt,
		UpdatedAt:       data.Data.UpdatedAt,
	}
}

func MapRefundsToTransactionModel(data *dto.PineOrderResponse, transactionID primitive.ObjectID) model.Transaction {
	// Map Payments
	payments := make([]model.Payment, 0)
	for _, p := range data.Data.Payments {
		payments = append(payments, model.Payment{
			ID:                       p.ID,
			MerchantPaymentReference: p.MerchantPaymentReference,
			Status:                   p.Status,
			PaymentMethod:            p.PaymentMethod,
			PaymentAmount: model.OrderAmount{
				Value:    p.PaymentAmount.Value,
				Currency: p.PaymentAmount.Currency,
			},
			AcquirerData: model.AcquirerData{
				ApprovalCode:      p.AcquirerData.ApprovalCode,
				AcquirerReference: p.AcquirerData.AcquirerReference,
				RRN:               p.AcquirerData.RRN,
				IsAggregator:      p.AcquirerData.IsAggregator,
				AcquirerName:      p.AcquirerData.AcquirerName,
			},
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}

	// Map Refunds
	refunds := make([]model.Refund, 0)
	for _, r := range data.Data.Refunds {
		// Map nested payments inside refunds
		refundPayments := make([]model.Payment, 0)
		for _, p := range r.Payments {
			refundPayments = append(refundPayments, model.Payment{
				ID:                       p.ID,
				MerchantPaymentReference: p.MerchantPaymentReference,
				Status:                   p.Status,
				PaymentMethod:            p.PaymentMethod,
				PaymentAmount: model.OrderAmount{
					Value:    p.PaymentAmount.Value,
					Currency: p.PaymentAmount.Currency,
				},
				AcquirerData: model.AcquirerData{
					ApprovalCode:      p.AcquirerData.ApprovalCode,
					AcquirerReference: p.AcquirerData.AcquirerReference,
					RRN:               p.AcquirerData.RRN,
					IsAggregator:      p.AcquirerData.IsAggregator,
					AcquirerName:      p.AcquirerData.AcquirerName,
				},
				CreatedAt: p.CreatedAt,
				UpdatedAt: p.UpdatedAt,
			})
		}

		refunds = append(refunds, model.Refund{
			TransactionID:          transactionID,
			MerchantOrderReference: r.MerchantOrderReference,
			OrderID:                r.OrderID,
			Type:                   r.Type,
			Status:                 r.Status,
			OrderAmount: model.OrderAmount{
				Value:    r.OrderAmount.Value,
				Currency: r.OrderAmount.Currency,
			},
			Payments: refundPayments,
			PurchaseDetails: model.PurchaseDetail{
				Customer: model.Customer{
					CountryCode:     r.PurchaseDetails.Customer.CountryCode,
					BillingAddress:  model.Address(r.PurchaseDetails.Customer.BillingAddress),
					ShippingAddress: model.Address(r.PurchaseDetails.Customer.ShippingAddress),
				},
				MerchantMetadata: model.MerchantMetadata(r.PurchaseDetails.MerchantMetadata),
			},
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}

	// Final Transaction Model with Refunds
	return model.Transaction{
		OrderId:                data.Data.OrderID,
		MerchantOrderReference: data.Data.MerchantOrderReference,
		OrderAmount: model.OrderAmount{
			Value:    data.Data.OrderAmount.Value,
			Currency: data.Data.OrderAmount.Currency,
		},
		PreAuth:               data.Data.PreAuth,
		AllowedPaymentMethods: data.Data.AllowedPaymentMethods,
		Notes:                 data.Data.Notes,
		Type:                  data.Data.Type,
		MerchantID:            data.Data.MerchantID,
		CallbackURL:           data.Data.CallbackURL,
		FailureCallbackURL:    data.Data.FailureCallbackURL,
		PurchaseDetails: model.PurchaseDetail{
			Customer: model.Customer{
				EmailID:         data.Data.PurchaseDetails.Customer.EmailID,
				FirstName:       data.Data.PurchaseDetails.Customer.FirstName,
				LastName:        data.Data.PurchaseDetails.Customer.LastName,
				CustomerID:      data.Data.PurchaseDetails.Customer.CustomerID,
				MobileNumber:    data.Data.PurchaseDetails.Customer.MobileNumber,
				BillingAddress:  model.Address(data.Data.PurchaseDetails.Customer.BillingAddress),
				ShippingAddress: model.Address(data.Data.PurchaseDetails.Customer.ShippingAddress),
			},
			MerchantMetadata: model.MerchantMetadata(data.Data.PurchaseDetails.MerchantMetadata),
		},
		PineOrderID:     data.Data.OrderID,
		Status:          data.Data.Status,
		Signature:       data.Data.Signature,
		IntegrationMode: data.Data.IntegrationMode,
		Payments:        payments,
		Refunds:         refunds, // ✅ Added refunds here
		CreatedAt:       data.Data.CreatedAt,
		UpdatedAt:       data.Data.UpdatedAt,
	}
}
