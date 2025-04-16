package services

import (
	"context"
	"fmt"
	"log"

	"github.com/aakritigkmit/payment-gateway/internal/dto"
	"github.com/aakritigkmit/payment-gateway/internal/model"
	"github.com/aakritigkmit/payment-gateway/internal/repository"
	"github.com/aakritigkmit/payment-gateway/internal/utils"
)

type ProductService struct {
	productRepo            *repository.ProductRepo
	productTransactionRepo *repository.ProductTransactionRepo
}

func NewProductService(productRepo *repository.ProductRepo, productTransactionRepo *repository.ProductTransactionRepo) *ProductService {
	return &ProductService{
		productRepo:            productRepo,
		productTransactionRepo: productTransactionRepo,
	}
}

func (s *ProductService) SyncProducts(ctx context.Context) error {
	products, err := utils.FetchDTOneProducts(ctx)
	if err != nil {
		return err
	}

	productChan := make(chan model.Product)
	errChan := make(chan error)
	done := make(chan bool)

	// Start 5 workers
	for i := 0; i < 5; i++ {
		go func() {
			for product := range productChan {
				if err := s.productRepo.FindOrCreateProduct(ctx, product); err != nil {
					errChan <- err
				}
			}
			done <- true
		}()
	}

	// Feed products to channel
	go func() {
		for _, product := range products {
			productChan <- product
		}
		close(productChan)
	}()

	// Wait for workers
	for i := 0; i < 5; i++ {
		<-done
	}
	close(errChan)

	for err := range errChan {
		log.Printf("failed to save product: %v", err)
	}

	return nil
}

func (s *ProductService) CreateAndSaveTransaction(ctx context.Context, req dto.CreateTransactionRequest) error {
	// Step 1: Create transaction via DT One
	if err := utils.CreateDTOneTransaction(ctx, req.ExternalID, req.ProductID, req.MobileNumber); err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// Step 2: Fetch transaction details
	productTransactions, err := utils.FetchDTOneTransactionByExternalID(ctx, req.ExternalID)
	if err != nil {
		return fmt.Errorf("failed to fetch transaction: %w", err)
	}

	// Step 3: Save to DB using ProductTransactionRepo
	for _, tx := range productTransactions {
		if err := s.productTransactionRepo.SaveProductTransaction(ctx, tx); err != nil {
			log.Printf("error saving product transaction: %v", err)
		}
	}

	return nil
}
