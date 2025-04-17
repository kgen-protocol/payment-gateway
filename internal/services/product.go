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
		wg       sync.WaitGroup
		sem      = make(chan struct{}, 3) // concurrency limit
		pinsChan = make(chan model.ProductPinItem, 1000)
		errChan  = make(chan error, 1000)
	)

	// Collector goroutine to gather all pins and errors
	var (
		allPins   []model.ProductPinItem
		allErrors []error
	)
	collectDone := make(chan struct{})
	go func() {
		for {
			select {
			case pin, ok := <-pinsChan:
				if !ok {
					pinsChan = nil
				} else {
					allPins = append(allPins, pin)
				}
			case err, ok := <-errChan:
				if !ok {
					errChan = nil
				} else {
					allErrors = append(allErrors, err)
				}
			}
			if pinsChan == nil && errChan == nil {
				break
			}
		}
		close(collectDone)
	}()

	for _, item := range req.LineItems {
		for i := 0; i < item.Quantity; i++ {
			wg.Add(1)
			sem <- struct{}{} // acquire slot

			go func(item dto.LineItem) {
				defer wg.Done()
				defer func() { <-sem }() // release slot

				externalID := fmt.Sprintf("TX-%s-%d", uuid.New().String()[:8], item.ProductID)

				if err := utils.CreateDTOneTransaction(ctx, externalID, item.ProductID, req.MobileNumber); err != nil {
					errChan <- fmt.Errorf("create failed for productID %d: %w", item.ProductID, err)
					return
				}

				var (
					transactions []model.ProductTransaction
					err          error
				)

				for attempt := 1; attempt <= 3; attempt++ {
					time.Sleep(time.Duration(1<<uint(attempt-1)) * 500 * time.Millisecond) // 500ms, 1s, 2s
					transactions, err = utils.FetchDTOneTransactionByExternalID(ctx, externalID)
					if err == nil && len(transactions) > 0 {
						break
					}
				}

				if err != nil {
					errChan <- fmt.Errorf("fetch failed for externalID %s: %w", externalID, err)
					return
				}

				for _, tx := range transactions {
					if err := s.productTransactionRepo.SaveProductTransaction(ctx, tx); err != nil {
						errChan <- fmt.Errorf("save failed for tx: %w", err)
						return
					}

					if tx.Pin.Code != "" && tx.Pin.Serial != "" {
						log.Printf("Adding PIN: %s / %s\n", tx.Pin.Code, tx.Pin.Serial)
						pinsChan <- model.ProductPinItem{
							ProductID: item.ProductID,
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
			}(item)
		}
	}

	wg.Wait()
	close(pinsChan)
	close(errChan)
	<-collectDone // wait for all pins/errors to be collected

	if len(allErrors) > 0 {
		return fmt.Errorf("failed to process bulk transactions: %w", allErrors[0])
	}

	log.Printf("Total product pins to save: %d\n", len(allPins))
	if len(allPins) > 0 {
		orderID := uuid.New().String()
		productPinData := model.ProductPin{
			OrderID:     orderID,
			ProductPins: allPins,
		}
		if err := s.productOrderRepo.SaveProductPins(ctx, productPinData); err != nil {
			return fmt.Errorf("failed to save product pins for orderID %s: %w", orderID, err)
		}
	}

	return nil
}
