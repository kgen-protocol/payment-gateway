package services

import (
	"context"
	"fmt"
	"time"

	"github.com/aakritigkmit/payment-gateway/internal/dto"
	"github.com/aakritigkmit/payment-gateway/internal/helpers"
	"github.com/aakritigkmit/payment-gateway/internal/model"
	"github.com/aakritigkmit/payment-gateway/internal/repository"
	"github.com/aakritigkmit/payment-gateway/internal/utils"
)

type OrderService struct {
	repo            *repository.OrderRepo
	transactionRepo *repository.TransactionRepo
}

func NewOrderService(repo *repository.OrderRepo, transactionRepo *repository.TransactionRepo) *OrderService {
	return &OrderService{repo, transactionRepo}
}

func (s *OrderService) PlaceOrder(ctx context.Context, req dto.PlaceOrderRequest) (utils.OrderAPIResponse, error) {
	// Fetch access token
	tokenResp, err := utils.FetchAccessToken(ctx)
	if err != nil {
		return utils.OrderAPIResponse{}, err
	}

	// Build JSON payload
	jsonPayload, err := helpers.BuildOrderPayload(req)
	if err != nil {
		return utils.OrderAPIResponse{}, err
	}

	// Create order on Pinelabs
	orderResp, err := utils.CreateOrderRequest(ctx, tokenResp.AccessToken, jsonPayload)
	if err != nil {
		return utils.OrderAPIResponse{}, err
	}

	// Prepare Order model for MongoDB
	order := &model.Order{
		OrderID:                orderResp.OrderID,
		MerchantOrderReference: string(rune(req.MerchantOrderReference)),
		Type:                   "ORDER",
		Status:                 "Pending",
		CallbackURL:            req.CallbackURL,
		FailureCallbackURL:     req.FailureCallbackURL,
		MerchantID:             "", // You can fill it from config or Pinelabs response if available
		Amount: model.OrderAmount{
			Value:    req.OrderAmount.Value,
			Currency: req.OrderAmount.Currency,
		},
		Notes:                 req.Notes,
		PreAuth:               req.PreAuth,
		AllowedPaymentMethods: req.AllowedPaymentMethods,
		IntegrationMode:       "REDIRECT", // Set based on config or your logic
		PurchaseDetails: model.PurchaseDetails{
			MerchantMetadata: req.PurchaseDetails.MerchantMetadata,
			Customer: model.Customer{
				EmailID:                      req.PurchaseDetails.Customer.EmailID,
				FirstName:                    req.PurchaseDetails.Customer.FirstName,
				LastName:                     req.PurchaseDetails.Customer.LastName,
				CustomerID:                   req.PurchaseDetails.Customer.CustomerID,
				MobileNumber:                 req.PurchaseDetails.Customer.MobileNumber,
				CountryCode:                  req.PurchaseDetails.Customer.CountryCode,
				IsEditCustomerDetailsAllowed: req.PurchaseDetails.Customer.IsEditCustomerDetailsAllowed,
				BillingAddress: model.Address{
					Address1: req.PurchaseDetails.Customer.BillingAddress.Address1,
					Address2: req.PurchaseDetails.Customer.BillingAddress.Address2,
					Address3: req.PurchaseDetails.Customer.BillingAddress.Address3,
					Pincode:  req.PurchaseDetails.Customer.BillingAddress.Pincode,
					City:     req.PurchaseDetails.Customer.BillingAddress.City,
					State:    req.PurchaseDetails.Customer.BillingAddress.State,
					Country:  req.PurchaseDetails.Customer.BillingAddress.Country,
				},
				ShippingAddress: model.Address{
					Address1: req.PurchaseDetails.Customer.ShippingAddress.Address1,
					Address2: req.PurchaseDetails.Customer.ShippingAddress.Address2,
					Address3: req.PurchaseDetails.Customer.ShippingAddress.Address3,
					Pincode:  req.PurchaseDetails.Customer.ShippingAddress.Pincode,
					City:     req.PurchaseDetails.Customer.ShippingAddress.City,
					State:    req.PurchaseDetails.Customer.ShippingAddress.State,
					Country:  req.PurchaseDetails.Customer.ShippingAddress.Country,
				},
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save Order to Mongo
	if err := s.repo.SaveOrder(ctx, order); err != nil {
		return utils.OrderAPIResponse{}, err
	}

	// No actual transaction exists yet (since order is just placed),
	// but if needed you can create a placeholder Transaction too.

	return utils.OrderAPIResponse{
		Token:       orderResp.Token,
		OrderID:     orderResp.OrderID,
		RedirectURL: orderResp.RedirectURL,
	}, nil
}

func (s *OrderService) SyncOrderDataFromPinelabs(ctx context.Context, pineOrderID string) error {

	// Fetch access token
	tokenResp, err := utils.FetchAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch access token: %w", err)
	}
	fmt.Println("-----------------------(1)")
	orderData, err := utils.GetOrderByIDFromPinelabs(ctx, tokenResp.AccessToken, pineOrderID)
	if err != nil {
		return fmt.Errorf("failed to fetch Pinelabs order: %w", err)
	}
	fmt.Println("-----------------------(2)")

	fmt.Println(orderData)

	// Map to MongoDB Order model
	order := &model.Order{
		OrderID:                orderData.OrderID,
		MerchantOrderReference: orderData.MerchantOrderReference,
		Type:                   orderData.Type,
		Status:                 orderData.Status,
		CallbackURL:            orderData.CallbackURL,
		FailureCallbackURL:     orderData.FailureCallbackURL,
		MerchantID:             orderData.MerchantID,
		Amount: model.OrderAmount{
			Value:    orderData.OrderAmount.Value,
			Currency: orderData.OrderAmount.Currency,
		},
		Notes:                   orderData.Notes,
		PreAuth:                 orderData.PreAuth,
		AllowedPaymentMethods:   orderData.AllowedPaymentMethods,
		CreatedAt:               orderData.CreatedAt,
		UpdatedAt:               orderData.UpdatedAt,
		IntegrationMode:         orderData.IntegrationMode,
		PaymentRetriesRemaining: orderData.PaymentRetriesRemaining,
		PurchaseDetails: model.PurchaseDetails{
			Customer: model.Customer{
				EmailID:                      orderData.PurchaseDetails.Customer.EmailID,
				FirstName:                    orderData.PurchaseDetails.Customer.FirstName,
				LastName:                     orderData.PurchaseDetails.Customer.LastName,
				CustomerID:                   orderData.PurchaseDetails.Customer.CustomerID,
				MobileNumber:                 orderData.PurchaseDetails.Customer.MobileNumber,
				CountryCode:                  orderData.PurchaseDetails.Customer.CountryCode,
				IsEditCustomerDetailsAllowed: orderData.PurchaseDetails.Customer.IsEditCustomerDetailsAllowed,
				BillingAddress:               model.Address(orderData.PurchaseDetails.Customer.BillingAddress),
				ShippingAddress:              model.Address(orderData.PurchaseDetails.Customer.ShippingAddress),
			},
			MerchantMetadata: orderData.PurchaseDetails.MerchantMetadata,
		},
	}

	// Save order
	if err := s.repo.UpdateOrder(ctx, order.OrderID, order); err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	// Convert and save each transaction/payment
	for _, payment := range orderData.Payments {
		transaction := &model.Transaction{
			TransactionID:            payment.ID,
			OrderID:                  orderData.OrderID,
			MerchantPaymentReference: payment.MerchantPaymentReference,
			Status:                   payment.Status,
			PaymentAmount: model.OrderAmount{
				Value:    payment.PaymentAmount.Value,
				Currency: payment.PaymentAmount.Currency,
			},
			PaymentMethod: payment.PaymentMethod,
			PaymentOption: model.PaymentOption{
				NetbankingData: func() *model.NetbankingData {
					if payment.PaymentOption.NetbankingData != nil {
						return &model.NetbankingData{
							PayCode: payment.PaymentOption.NetbankingData.PayCode,
							TxnMode: payment.PaymentOption.NetbankingData.TxnMode,
						}
					}
					return nil
				}(),
				UPIData: func() *model.UPIData {
					if payment.PaymentOption.UPIData != nil {
						return &model.UPIData{
							TxnMode: payment.PaymentOption.UPIData.TxnMode,
						}
					}
					return nil
				}(),
			},
			AcquirerData: model.AcquirerData{
				ApprovalCode:      payment.AcquirerData.ApprovalCode,
				AcquirerReference: payment.AcquirerData.AcquirerReference,
				RRN:               payment.AcquirerData.RRN,
				IsAggregator:      payment.AcquirerData.IsAggregator,
			},
			ErrorDetail: model.ErrorDetail{
				Code:    payment.ErrorDetail.Code,
				Message: payment.ErrorDetail.Message,
			},
			CreatedAt: payment.CreatedAt,
			UpdatedAt: payment.UpdatedAt,
		}

		if err := s.transactionRepo.SaveTransaction(ctx, transaction); err != nil {
			return fmt.Errorf("failed to save transaction: %w", err)
		}
	}

	return nil
}
