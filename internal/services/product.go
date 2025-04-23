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

// func (s *ProductService) CreateAndSaveBulkTransactions(ctx context.Context, req dto.BulkTransactionRequest) error {
// 	startTime := time.Now()
// 	const (
// 		numWorkers = 10
// 		maxRetries = 5
// 		baseDelay  = 2 * time.Second
// 	)

// 	type task struct {
// 		LineItem   dto.LineItem
// 		ExternalID string
// 	}

// 	// Total quantity from request
// 	expectedQty := 0
// 	for _, item := range req.LineItems {
// 		expectedQty += item.Quantity
// 	}

// 	taskChan := make(chan task)
// 	resultChan := make(chan model.ProductPinItem, expectedQty) // buffered to avoid blocking
// 	errorChan := make(chan error, expectedQty)

// 	var wg sync.WaitGroup

// 	// Start workers
// 	for w := 0; w < numWorkers; w++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			for t := range taskChan {
// 				txRecord := model.ProductTransaction{
// 					ExternalID: t.ExternalID,
// 					CreatedAt:  time.Now(),
// 					UpdatedAt:  time.Now(),
// 				}

// 				// Save initial transaction
// 				if err := s.productTransactionRepo.SaveProductTransaction(ctx, txRecord); err != nil {
// 					errorChan <- fmt.Errorf("initial save failed for %s: %w", t.ExternalID, err)
// 					continue
// 				}

// 				// Create DT One transaction
// 				var err error
// 				for attempt := 1; attempt <= maxRetries; attempt++ {
// 					err = utils.CreateDTOneTransaction(ctx, t.ExternalID, t.LineItem.ProductID, req.MobileNumber)
// 					if err == nil {
// 						break
// 					}
// 					time.Sleep(baseDelay * time.Duration(attempt))
// 				}
// 				if err != nil {
// 					errorChan <- fmt.Errorf("CreateTX failed for %s: %w", t.ExternalID, err)
// 					continue
// 				}

// 				// Fetch transaction details
// 				var txs []model.ProductTransaction
// 				for attempt := 1; attempt <= maxRetries; attempt++ {
// 					time.Sleep(baseDelay * time.Duration(attempt))
// 					txs, err = utils.FetchDTOneTransactionByExternalID(ctx, t.ExternalID)
// 					if err == nil && len(txs) > 0 {
// 						break
// 					}
// 				}
// 				if err != nil || len(txs) == 0 {
// 					errorChan <- fmt.Errorf("FetchTX failed for %s: %w", t.ExternalID, err)
// 					continue
// 				}

// 				// Update and collect pin
// 				for _, tx := range txs {
// 					tx.UpdatedAt = time.Now()
// 					if err := s.productTransactionRepo.UpdateProductTransaction(ctx, tx.ExternalID, tx); err != nil {
// 						errorChan <- fmt.Errorf("UpdateTX failed for %s: %w", tx.ExternalID, err)
// 						continue
// 					}

// 					resultChan <- model.ProductPinItem{
// 						ExternalID: tx.ExternalID,
// 						ProductID:  t.LineItem.ProductID,
// 						Pin: struct {
// 							Code   string `bson:"code"`
// 							Serial string `bson:"serial"`
// 						}{
// 							Code:   tx.Pin.Code,
// 							Serial: tx.Pin.Serial,
// 						},
// 					}
// 				}
// 			}
// 		}()
// 	}

// 	// Collector
// 	var (
// 		successPins = make([]model.ProductPinItem, 0, expectedQty)
// 		saveErrors  = make([]error, 0)
// 	)
// 	collectDone := make(chan struct{})
// 	go func() {
// 		for i := 0; i < expectedQty; i++ {
// 			select {
// 			case pin := <-resultChan:
// 				successPins = append(successPins, pin)
// 			case err := <-errorChan:
// 				saveErrors = append(saveErrors, err)
// 			}
// 		}
// 		close(collectDone)
// 	}()

// 	// Push tasks
// 	go func() {
// 		for _, item := range req.LineItems {
// 			for i := 0; i < item.Quantity; i++ {
// 				taskChan <- task{
// 					LineItem:   item,
// 					ExternalID: fmt.Sprintf("TX-%s-%d", uuid.New().String()[:8], item.ProductID),
// 				}
// 			}
// 		}
// 		close(taskChan)
// 	}()

// 	// Wait for workers and collector
// 	wg.Wait()
// 	close(resultChan)
// 	close(errorChan)
// 	<-collectDone

// 	if len(successPins) > 0 {
// 		orderID := uuid.New().String()
// 		pinDoc := model.ProductPin{
// 			OrderID:     orderID,
// 			ProductPins: successPins,
// 			CreatedAt:   time.Now(),
// 			UpdatedAt:   time.Now(),
// 		}
// 		if err := s.productOrderRepo.SaveProductPins(ctx, pinDoc); err != nil {
// 			log.Printf("SaveProductPins error: %v", err)
// 			saveErrors = append(saveErrors, err)
// 		} else {
// 			log.Printf("Saved %d pins under OrderID %s", len(successPins), orderID)
// 		}
// 	}

// 	// Final log
// 	log.Printf("Reconciliation: Expected=%d, Saved=%d, Errors=%d", expectedQty, len(successPins), len(saveErrors))
// 	log.Printf("Total execution time: %v", time.Since(startTime)) // <-- Add before success return
// 	if len(saveErrors) > 0 {
// 		for i, err := range saveErrors {
// 			fmt.Printf("[%d] %v\n", i+1, err)
// 		}
// 		return fmt.Errorf("completed with %d errors", len(saveErrors))
// 	}

// 	return nil
// }

func (s *ProductService) InitBulkProductTransaction(ctx context.Context, req dto.BulkTransactionRequest) (string, error) {
	orderId := uuid.New().String()

	productOrder := model.ProductPin{
		OrderID:     orderId,
		ProductPins: []model.ProductPinItem{}, // initially empty, to be filled later
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.productOrderRepo.SaveProductPins(ctx, productOrder); err != nil {
		return "", err
	}
	return orderId, nil
}

func (s *ProductService) ProcessBulkProductTransactionAsync(ctx context.Context, req dto.BulkTransactionRequest, orderId string) error {
	startTime := time.Now()

	log.Printf("[INFO] Started processing bulk transaction for OrderID: %s", orderId)

	type task struct {
		LineItem   dto.LineItem
		ExternalID string
	}

	const numWorkers = 10

	var (
		wg             sync.WaitGroup
		taskChan       = make(chan task)
		resultChan     = make(chan model.ProductPinItem)
		retryTasks     = make([]task, 0)
		retryFetchOnly = make([]task, 0)
		mu             sync.Mutex
	)

	// Workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			log.Printf("[INFO] Worker %d started", workerID)

			for t := range taskChan {
				log.Printf("[DEBUG] [Worker %d] Handling task ExternalID: %s, ProductID: %d", workerID, t.ExternalID, t.LineItem.ProductID)

				txRecord := model.ProductTransaction{
					ExternalID: t.ExternalID,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}
				if err := s.productTransactionRepo.SaveProductTransaction(ctx, txRecord); err != nil {
					log.Printf("[ERROR] [Worker %d] Initial save failed for %s: %v", workerID, t.ExternalID, err)
					continue
				}

				if err := utils.CreateDTOneTransaction(ctx, t.ExternalID, t.LineItem.ProductID, req.MobileNumber); err != nil {
					log.Printf("[WARN] [Worker %d] CreateTX failed for %s: %v", workerID, t.ExternalID, err)
					mu.Lock()
					retryTasks = append(retryTasks, t)
					mu.Unlock()
					continue
				}

				time.Sleep(300 * time.Millisecond) // avoid hammering API

				txs, err := utils.FetchDTOneTransactionByExternalID(ctx, t.ExternalID)
				if err != nil || len(txs) == 0 {
					log.Printf("[WARN] [Worker %d] FetchTX failed for %s: %v", workerID, t.ExternalID, err)
					mu.Lock()
					retryFetchOnly = append(retryFetchOnly, t)
					mu.Unlock()
					continue
				}

				for _, tx := range txs {
					tx.UpdatedAt = time.Now()
					if err := s.productTransactionRepo.UpdateProductTransaction(ctx, tx.ExternalID, tx); err != nil {
						log.Printf("[ERROR] [Worker %d] UpdateTX failed for %s: %v", workerID, tx.ExternalID, err)
						continue
					}
					log.Printf("[INFO] [Worker %d] TX fetched - ExternalID: %s, ProductID: %d, Pin: %s, Serial: %s",
						workerID, tx.ExternalID, tx.Product.UniqueId, tx.Pin.Code, tx.Pin.Serial)

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
			log.Printf("[INFO] Worker %d exited", workerID)
		}(i)
	}

	go func() {
		log.Printf("[DEBUG] Dispatching tasks...")
		for _, item := range req.LineItems {
			for i := 0; i < item.Quantity; i++ {
				externalID := fmt.Sprintf("TX-%s-%d", uuid.New().String()[:8], item.ProductID)
				log.Printf("[DEBUG] Queuing task for ProductID: %d, ExternalID: %s", item.ProductID, externalID)
				taskChan <- task{LineItem: item, ExternalID: externalID}
			}
		}
		close(taskChan)
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var pinItems []model.ProductPinItem
	for pin := range resultChan {
		pinItems = append(pinItems, pin)
	}

	// Retry failed transactions
	for _, t := range retryTasks {
		log.Printf("[RETRY] Retrying CreateTX for ExternalID: %s", t.ExternalID)
		if err := utils.CreateDTOneTransaction(ctx, t.ExternalID, t.LineItem.ProductID, req.MobileNumber); err != nil {
			log.Printf("[ERROR] Final CreateTX failed for %s: %v", t.ExternalID, err)
			continue
		}
		retryFetchOnly = append(retryFetchOnly, t)
	}

	for _, t := range retryFetchOnly {
		time.Sleep(300 * time.Millisecond)
		log.Printf("[RETRY] Fetching transaction post-retry for ExternalID: %s", t.ExternalID)

		txs, err := utils.FetchDTOneTransactionByExternalID(ctx, t.ExternalID)
		if err != nil || len(txs) == 0 {
			log.Printf("[ERROR] Final FetchTX failed for %s: %v", t.ExternalID, err)
			continue
		}

		for _, tx := range txs {
			tx.UpdatedAt = time.Now()
			if err := s.productTransactionRepo.UpdateProductTransaction(ctx, tx.ExternalID, tx); err != nil {
				log.Printf("[ERROR] UpdateTX (after retry) failed for %s: %v", tx.ExternalID, err)
				continue
			}

			pinItems = append(pinItems, model.ProductPinItem{
				ExternalID: tx.ExternalID,
				ProductID:  t.LineItem.ProductID,
				Pin: struct {
					Code   string `bson:"code"`
					Serial string `bson:"serial"`
				}{
					Code:   tx.Pin.Code,
					Serial: tx.Pin.Serial,
				},
			})
		}
	}

	// Dump pins
	if len(pinItems) > 0 {
		log.Printf("[INFO] Dumping %d pins for OrderID: %s", len(pinItems), orderId)
		dump := model.ProductPinDump{
			OrderID:     orderId,
			ProductPins: pinItems,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := s.productOrderRepo.SaveProductPinsDump(ctx, dump); err != nil {
			log.Printf("[ERROR] Error saving pin dump for OrderID %s: %v", orderId, err)
		}
	}

	// Final update to ProductOrder
	finalPins, err := s.productOrderRepo.GetPinsByOrderID(ctx, orderId)
	if err != nil {
		log.Printf("[ERROR] Error fetching dumped pins for OrderID %s: %v", orderId, err)
		return err
	}

	log.Printf("Total execution time for OrderID %s: %v", orderId, time.Since(startTime))

	if err := s.productOrderRepo.UpdateProductOrderWithPins(ctx, orderId, finalPins); err != nil {
		log.Printf("[ERROR] Error updating ProductOrder with final pins for OrderID %s: %v", orderId, err)
		return err
	}

	log.Printf("[SUCCESS] OrderID %s processed with %d total pins", orderId, len(finalPins))
	return nil
}

// func (s *ProductService) ProcessBulkProductTransactionAsync(ctx context.Context, req dto.BulkTransactionRequest, orderId string) error {
// 	startTime := time.Now()

// 	const (
// 		numWorkers = 10
// 		maxRetries = 5
// 		baseDelay  = 2 * time.Second
// 	)

// 	type task struct {
// 		LineItem   dto.LineItem
// 		ExternalID string
// 		OrderID    string
// 	}

// 	expectedQty := 0
// 	for _, item := range req.LineItems {
// 		expectedQty += item.Quantity
// 	}

// 	log.Printf("[OrderID %s] Starting async processing with %d expected items", orderId, expectedQty)

// 	taskChan := make(chan task)
// 	resultChan := make(chan model.ProductPinItem, expectedQty)
// 	errorChan := make(chan error, expectedQty)

// 	var wg sync.WaitGroup

// 	for w := 0; w < numWorkers; w++ {
// 		wg.Add(1)
// 		go func(workerID int) {
// 			defer wg.Done()
// 			log.Printf("[Worker %d] Started", workerID)
// 			for t := range taskChan {
// 				log.Printf("[Worker %d] Processing task: ExternalID=%s, ProductID=%d", workerID, t.ExternalID, t.LineItem.ProductID)

// 				txRecord := model.ProductTransaction{
// 					ExternalID: t.ExternalID,
// 					CreatedAt:  time.Now(),
// 					UpdatedAt:  time.Now(),
// 				}

// 				if err := s.productTransactionRepo.SaveProductTransaction(ctx, txRecord); err != nil {
// 					errorChan <- fmt.Errorf("initial save failed for %s: %w", t.ExternalID, err)
// 					continue
// 				}
// 				log.Printf("[Worker %d] Saved initial transaction for %s", workerID, t.ExternalID)

// 				var err error
// 				for attempt := 1; attempt <= maxRetries; attempt++ {
// 					err = utils.CreateDTOneTransaction(ctx, t.ExternalID, t.LineItem.ProductID, req.MobileNumber)
// 					if err == nil {
// 						break
// 					}
// 					log.Printf("[Worker %d] CreateTX retry %d for %s failed: %v", workerID, attempt, t.ExternalID, err)
// 					time.Sleep(baseDelay * time.Duration(attempt))
// 				}
// 				if err != nil {
// 					errorChan <- fmt.Errorf("CreateTX failed for %s: %w", t.ExternalID, err)
// 					continue
// 				}
// 				log.Printf("[Worker %d] CreateTX successful for %s", workerID, t.ExternalID)

// 				var txs []model.ProductTransaction
// 				for attempt := 1; attempt <= maxRetries; attempt++ {
// 					time.Sleep(baseDelay * time.Duration(attempt))
// 					txs, err = utils.FetchDTOneTransactionByExternalID(ctx, t.ExternalID)
// 					if err == nil && len(txs) > 0 {
// 						break
// 					}
// 					log.Printf("[Worker %d] FetchTX retry %d for %s: %v", workerID, attempt, t.ExternalID, err)
// 				}
// 				if err != nil || len(txs) == 0 {
// 					errorChan <- fmt.Errorf("FetchTX failed for %s: %w", t.ExternalID, err)
// 					continue
// 				}
// 				log.Printf("[Worker %d] FetchTX successful for %s", workerID, t.ExternalID)

// 				for _, tx := range txs {
// 					tx.UpdatedAt = time.Now()
// 					if err := s.productTransactionRepo.UpdateProductTransaction(ctx, tx.ExternalID, tx); err != nil {
// 						errorChan <- fmt.Errorf("UpdateTX failed for %s: %w", tx.ExternalID, err)
// 						continue
// 					}
// 					log.Printf("[Worker %d] UpdateTX successful for %s", workerID, tx.ExternalID)

// 					resultChan <- model.ProductPinItem{
// 						ExternalID: tx.ExternalID,
// 						ProductID:  t.LineItem.ProductID,
// 						Pin: struct {
// 							Code   string `bson:"code"`
// 							Serial string `bson:"serial"`
// 						}{
// 							Code:   tx.Pin.Code,
// 							Serial: tx.Pin.Serial,
// 						},
// 					}
// 				}
// 			}
// 			log.Printf("[Worker %d] Finished", workerID)
// 		}(w)
// 	}

// 	successPins := make([]model.ProductPinItem, 0, expectedQty)
// 	saveErrors := make([]error, 0)
// 	collectDone := make(chan struct{})

// 	go func() {
// 		for i := 0; i < expectedQty; i++ {
// 			select {
// 			case pin := <-resultChan:
// 				log.Printf("Collected PIN: ExternalID=%s, ProductID=%d", pin.ExternalID, pin.ProductID)
// 				successPins = append(successPins, pin)
// 			case err := <-errorChan:
// 				log.Printf("Error collected: %v", err)
// 				saveErrors = append(saveErrors, err)
// 			}
// 		}
// 		close(collectDone)
// 	}()

// 	go func() {
// 		for _, item := range req.LineItems {
// 			for i := 0; i < item.Quantity; i++ {
// 				externalID := fmt.Sprintf("TX-%s-%d", uuid.New().String()[:8], item.ProductID)
// 				taskChan <- task{
// 					LineItem:   item,
// 					ExternalID: externalID,
// 					OrderID:    orderId,
// 				}
// 				log.Printf("Dispatched task: ExternalID=%s, ProductID=%d", externalID, item.ProductID)
// 			}
// 		}
// 		close(taskChan)
// 	}()

// 	wg.Wait()
// 	close(resultChan)
// 	close(errorChan)
// 	<-collectDone

// 	log.Printf("[OrderID %s] All workers and collector finished", orderId)

// 	if len(successPins) > 0 {
// 		if err := s.productOrderRepo.UpdateProductPins(ctx, orderId, successPins); err != nil {
// 			log.Printf("UpdateProductPins error for OrderID %s: %v", orderId, err)
// 			saveErrors = append(saveErrors, err)
// 		} else {
// 			log.Printf("Successfully updated %d pins for OrderID %s", len(successPins), orderId)
// 		}
// 	}

// 	log.Printf("Reconciliation for OrderID %s: Expected=%d, Saved=%d, Errors=%d", orderId, expectedQty, len(successPins), len(saveErrors))
// 	log.Printf("Total execution time for OrderID %s: %v", orderId, time.Since(startTime))

// 	if len(saveErrors) > 0 {
// 		for i, err := range saveErrors {
// 			log.Printf("[Error %d] %v", i+1, err)
// 		}
// 		return fmt.Errorf("completed with %d errors", len(saveErrors))
// 	}

// 	return nil
// }
