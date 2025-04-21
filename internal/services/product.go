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
	startTime := time.Now()
	const (
		numWorkers = 10
		maxRetries = 5
		baseDelay  = 2 * time.Second
	)

	type task struct {
		LineItem   dto.LineItem
		ExternalID string
	}

	// Total quantity from request
	expectedQty := 0
	for _, item := range req.LineItems {
		expectedQty += item.Quantity
	}

	taskChan := make(chan task)
	resultChan := make(chan model.ProductPinItem, expectedQty) // buffered to avoid blocking
	errorChan := make(chan error, expectedQty)

	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range taskChan {
				txRecord := model.ProductTransaction{
					ExternalID: t.ExternalID,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}

				// Save initial transaction
				if err := s.productTransactionRepo.SaveProductTransaction(ctx, txRecord); err != nil {
					errorChan <- fmt.Errorf("initial save failed for %s: %w", t.ExternalID, err)
					continue
				}

				// Create DT One transaction
				var err error
				for attempt := 1; attempt <= maxRetries; attempt++ {
					err = utils.CreateDTOneTransaction(ctx, t.ExternalID, t.LineItem.ProductID, req.MobileNumber)
					if err == nil {
						break
					}
					time.Sleep(baseDelay * time.Duration(attempt))
				}
				if err != nil {
					errorChan <- fmt.Errorf("CreateTX failed for %s: %w", t.ExternalID, err)
					continue
				}

				// Fetch transaction details
				var txs []model.ProductTransaction
				for attempt := 1; attempt <= maxRetries; attempt++ {
					time.Sleep(baseDelay * time.Duration(attempt))
					txs, err = utils.FetchDTOneTransactionByExternalID(ctx, t.ExternalID)
					if err == nil && len(txs) > 0 {
						break
					}
				}
				if err != nil || len(txs) == 0 {
					errorChan <- fmt.Errorf("FetchTX failed for %s: %w", t.ExternalID, err)
					continue
				}

				// Update and collect pin
				for _, tx := range txs {
					tx.UpdatedAt = time.Now()
					if err := s.productTransactionRepo.UpdateProductTransaction(ctx, tx.ExternalID, tx); err != nil {
						errorChan <- fmt.Errorf("UpdateTX failed for %s: %w", tx.ExternalID, err)
						continue
					}

					resultChan <- model.ProductPinItem{
						ExternalID: tx.ExternalID,
						ProductID:  t.LineItem.ProductID,
						Pin: struct {
							Code   string `bson:"code"`
							Serial string `bson:"serial"`
						}{
							Code:   tx.Pin.Code,
							Serial: tx.Pin.Serial,
						},
					}
				}
			}
		}()
	}

	// Collector
	var (
		successPins = make([]model.ProductPinItem, 0, expectedQty)
		saveErrors  = make([]error, 0)
	)
	collectDone := make(chan struct{})
	go func() {
		for i := 0; i < expectedQty; i++ {
			select {
			case pin := <-resultChan:
				successPins = append(successPins, pin)
			case err := <-errorChan:
				saveErrors = append(saveErrors, err)
			}
		}
		close(collectDone)
	}()

	// Push tasks
	go func() {
		for _, item := range req.LineItems {
			for i := 0; i < item.Quantity; i++ {
				taskChan <- task{
					LineItem:   item,
					ExternalID: fmt.Sprintf("TX-%s-%d", uuid.New().String()[:8], item.ProductID),
				}
			}
		}
		close(taskChan)
	}()

	// Wait for workers and collector
	wg.Wait()
	close(resultChan)
	close(errorChan)
	<-collectDone

	if len(successPins) > 0 {
		orderID := uuid.New().String()
		pinDoc := model.ProductPin{
			OrderID:     orderID,
			ProductPins: successPins,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := s.productOrderRepo.SaveProductPins(ctx, pinDoc); err != nil {
			log.Printf("SaveProductPins error: %v", err)
			saveErrors = append(saveErrors, err)
		} else {
			log.Printf("Saved %d pins under OrderID %s", len(successPins), orderID)
		}
	}

	// Final log
	log.Printf("Reconciliation: Expected=%d, Saved=%d, Errors=%d", expectedQty, len(successPins), len(saveErrors))
	log.Printf("Total execution time: %v", time.Since(startTime)) // <-- Add before success return
	if len(saveErrors) > 0 {
		for i, err := range saveErrors {
			fmt.Printf("[%d] %v\n", i+1, err)
		}
		return fmt.Errorf("completed with %d errors", len(saveErrors))
	}

	return nil
}
