package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aakritigkmit/payment-gateway/internal/dto"
	"github.com/aakritigkmit/payment-gateway/internal/model"
	"github.com/aakritigkmit/payment-gateway/internal/repository"
	"github.com/aakritigkmit/payment-gateway/internal/utils"
	"github.com/google/uuid"
)

type ProductService struct {
	productRepo            *repository.ProductRepo
	productTransactionRepo *repository.ProductTransactionRepo
	productOrderRepo       *repository.ProductOrderRepo
}

func NewProductService(productRepo *repository.ProductRepo, productTransactionRepo *repository.ProductTransactionRepo, productOrderRepo *repository.ProductOrderRepo) *ProductService {
	return &ProductService{
		productRepo:            productRepo,
		productTransactionRepo: productTransactionRepo,
		productOrderRepo:       productOrderRepo,
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

func (s *ProductService) CreateAndSaveBulkTransactions(ctx context.Context, req dto.BulkTransactionRequest) error {
	var (
		wg         sync.WaitGroup
		sem        = make(chan struct{}, 10) // max 10 concurrent transactions
		mu         sync.Mutex
		totalPins  []model.ProductPinItem
		saveErrors []error
	)

	for _, item := range req.LineItems {
		for i := 0; i < item.Quantity; i++ {
			wg.Add(1)
			sem <- struct{}{}

			go func(item dto.LineItem) {
				defer wg.Done()
				defer func() { <-sem }()

				externalID := fmt.Sprintf("TX-%s-%d", uuid.New().String()[:8], item.ProductID)

				// Step 1: Create transaction
				err := utils.CreateDTOneTransaction(ctx, externalID, item.ProductID, req.MobileNumber)
				if err != nil {
					log.Printf("Error creating transaction for ProductID %d: %v", item.ProductID, err)
					mu.Lock()
					saveErrors = append(saveErrors, err)
					mu.Unlock()
					return
				}

				// Step 2: Fetch transaction with retries
				var txs []model.ProductTransaction
				for attempt := 1; attempt <= 5; attempt++ {
					time.Sleep(time.Duration(attempt) * time.Second) // increasing delay
					txs, err = utils.FetchDTOneTransactionByExternalID(ctx, externalID)
					if err == nil && len(txs) > 0 {
						break
					}
				}

				if err != nil || len(txs) == 0 {
					log.Printf("Failed to fetch transaction for externalID %s: %v", externalID, err)
					mu.Lock()
					saveErrors = append(saveErrors, fmt.Errorf("fetch failed for %s", externalID))
					mu.Unlock()
					return
				}

				// Step 3: Save transaction and extract PIN
				for _, tx := range txs {
					tx.CreatedAt = time.Now()
					tx.UpdatedAt = time.Now()
					if err := s.productTransactionRepo.SaveProductTransaction(ctx, tx); err != nil {
						log.Printf("Error saving transaction: %v", err)
						mu.Lock()
						saveErrors = append(saveErrors, err)
						mu.Unlock()
						continue
					}

					if tx.Pin.Code != "" && tx.Pin.Serial != "" {
						pin := model.ProductPinItem{
							ProductID: item.ProductID,
							Pin: struct {
								Code   string `bson:"code"`
								Serial string `bson:"serial"`
							}{
								Code:   tx.Pin.Code,
								Serial: tx.Pin.Serial,
							},
						}
						mu.Lock()
						totalPins = append(totalPins, pin)
						mu.Unlock()
					}
				}
			}(item)
		}
	}

	wg.Wait()

	// Reconciliation check
	expectedPinCount := 0
	for _, item := range req.LineItems {
		expectedPinCount += item.Quantity
	}

	if len(totalPins) != expectedPinCount {
		log.Printf("WARNING: Expected %d PINs, but got %d", expectedPinCount, len(totalPins))
	}

	if len(totalPins) > 0 {
		orderID := uuid.New().String()
		productPinData := model.ProductPin{
			OrderID:     orderID,
			ProductPins: totalPins,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := s.productOrderRepo.SaveProductPins(ctx, productPinData); err != nil {
			log.Printf("Failed to save final product pins: %v", err)
			return err
		}
		log.Printf("Saved %d PINs successfully to order %s", len(totalPins), orderID)
	}

	if len(saveErrors) > 0 {
		return fmt.Errorf("completed with %d errors, see logs for details", len(saveErrors))
	}

	return nil
}
