package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/aakritigkmit/payment-gateway/internal/dto"
	"github.com/aakritigkmit/payment-gateway/internal/model"
	"github.com/aakritigkmit/payment-gateway/internal/repository"
	"github.com/aakritigkmit/payment-gateway/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
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

func (s *ProductService) SyncProducts(ctx context.Context, filter dto.ProductSyncRequest) error {
	startTime := time.Now()

	const perPage = 100
	const fetchConcurrency = 10
	const dbWorkerCount = 20

	log.Println("Starting DT One product sync...")

	_, totalPages, err := utils.FetchDTOneProducts(ctx, 1, perPage, filter)
	if err != nil {
		log.Printf("Failed to fetch initial product page: %v", err)
		return err
	}
	log.Printf("Total pages to sync: %d", totalPages)

	productChan := make(chan model.Product, 2000)
	saveErrChan := make(chan ProductSaveError, 1000)
	fetchErrChan := make(chan FetchPageError, 1000)
	successChan := make(chan int, 2000)

	var dbWg sync.WaitGroup
	dbWg.Add(dbWorkerCount)

	for i := 0; i < dbWorkerCount; i++ {
		go func(workerID int) {
			defer dbWg.Done()
			for product := range productChan {
				if err := s.productRepo.FindOrCreateProduct(ctx, product); err != nil {
					saveErrChan <- ProductSaveError{
						ProductID: product.UniqueId,
						ErrorMsg:  fmt.Sprintf("worker %d: failed to save product: %v", workerID, err),
					}
				} else {
					successChan <- product.UniqueId
				}
			}
		}(i)
	}

	var fetchWg sync.WaitGroup
	pageChan := make(chan int, totalPages)

	go func() {
		for p := 1; p <= totalPages; p++ {
			pageChan <- p
		}
		close(pageChan)
	}()

	for i := 0; i < fetchConcurrency; i++ {
		fetchWg.Add(1)
		go func(workerID int) {
			defer fetchWg.Done()
			for page := range pageChan {
				select {
				case <-ctx.Done():
					return
				default:
					log.Printf("[Fetcher %d] Fetching page %d...", workerID, page)
					pageProducts, _, err := utils.FetchDTOneProducts(ctx, page, perPage, filter)
					if err != nil {
						fetchErrChan <- FetchPageError{
							Page:     page,
							ErrorMsg: fmt.Sprintf("fetcher %d: page %d fetch error: %v", workerID, page, err),
						}
						continue
					}
					log.Printf("[Fetcher %d] Page %d: fetched %d products", workerID, page, len(pageProducts))
					for _, product := range pageProducts {
						productChan <- product
					}
				}
			}
		}(i)
	}

	// Closing channels after fetchers and db workers complete
	go func() {
		fetchWg.Wait()
		close(productChan)
	}()

	go func() {
		dbWg.Wait()
		close(successChan)
		close(saveErrChan)
		close(fetchErrChan)
	}()

	// Collect results
	var (
		successIDs        []int
		productSaveErrors []ProductSaveError
		fetchPageErrors   []FetchPageError
	)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		for id := range successChan {
			successIDs = append(successIDs, id)
		}
	}()

	go func() {
		defer wg.Done()
		for e := range saveErrChan {
			log.Printf("Save Error: %v", e)
			productSaveErrors = append(productSaveErrors, e)
		}
	}()

	go func() {
		defer wg.Done()
		for e := range fetchErrChan {
			log.Printf("Fetch Error: %v", e)
			fetchPageErrors = append(fetchPageErrors, e)
		}
	}()

	wg.Wait()

	log.Printf("DT One product sync complete. Success: %d, Save Errors: %d, Fetch Errors: %d.", len(successIDs), len(productSaveErrors), len(fetchPageErrors))
	log.Printf("Total execution time: %v", time.Since(startTime))

	// Generate the XLSX Report
	err = GenerateProductSyncReport(len(successIDs), len(productSaveErrors), len(fetchPageErrors), productSaveErrors, fetchPageErrors)
	if err != nil {
		log.Printf("Failed to generate sync report: %v", err)
	}

	return nil
}

func (s *ProductService) GenerateProductFetchReport(ctx context.Context, filter dto.ProductSyncRequest) error {
	startTime := time.Now()

	const perPage = 100
	const fetchConcurrency = 10

	start := time.Now()
	log.Println("[Report] Starting DT One product report generation...")

	// Step 1: Initial fetch to get total pages
	_, totalPages, err := utils.FetchDTOneProducts(ctx, 1, perPage, filter)
	if err != nil {
		log.Printf("[Report] Initial fetch failed: %v", err)
		return err
	}
	log.Printf("[Report] Total pages to fetch: %d", totalPages)

	productChan := make(chan model.Product, 5000)
	pageChan := make(chan int, totalPages)

	// Step 2: Producer - feed page numbers
	go func() {
		for i := 1; i <= totalPages; i++ {
			log.Printf("[Producer] Queuing page %d", i)
			pageChan <- i
		}
		close(pageChan)
		log.Println("[Producer] All pages queued.")
	}()

	// Step 3: Fetchers - concurrent data fetch
	var fetchWg sync.WaitGroup
	for i := 0; i < fetchConcurrency; i++ {
		fetchWg.Add(1)
		go func(workerID int) {
			defer fetchWg.Done()
			log.Printf("[Fetcher %d] Started.", workerID)
			for page := range pageChan {
				select {
				case <-ctx.Done():
					log.Printf("[Fetcher %d] Context cancelled. Exiting.", workerID)
					return
				default:
					products, _, err := utils.FetchDTOneProducts(ctx, page, perPage, filter)
					if err != nil {
						log.Printf("[Fetcher %d] Page %d error: %v", workerID, page, err)
						continue
					}
					log.Printf("[Fetcher %d] Fetched %d products from page %d", workerID, len(products), page)
					for _, p := range products {
						productChan <- p
					}
				}
			}
			log.Printf("[Fetcher %d] Completed.", workerID)
		}(i)
	}

	// Step 4: Closer - close product channel after all fetchers complete
	go func() {
		fetchWg.Wait()
		close(productChan)
		log.Println("[Collector] All fetchers done. Product channel closed.")
	}()

	// Step 5: Collector - gather all products
	var allProducts []model.Product
	count := 0
	for product := range productChan {
		allProducts = append(allProducts, product)
		count++
		if count%100 == 0 {
			log.Printf("[Collector] Collected %d products so far...", count)
		}
	}

	log.Printf("[Report] Fetched total %d products in %v", len(allProducts), time.Since(start))

	// Step 6: Generate Excel report
	log.Println("[Report] Generating Excel report...")
	err = ExportProductsToExcel(allProducts)
	if err != nil {
		log.Printf("[Report] Excel generation failed: %v", err)
		return err
	}
	log.Println("[Report] Excel report generated successfully.")
	log.Printf("Total execution time: %v", time.Since(startTime))

	return nil
}

func (s *ProductService) GenerateProductReportByIDs(ctx context.Context, productIDs []int) error {
	log.Println("[ReportByIDs] Starting product report generation by IDs...")
	start := time.Now()

	products, err := s.productRepo.GetProductsByUniqueIDs(ctx, productIDs)
	if err != nil {
		log.Printf("[ReportByIDs] DB fetch failed: %v", err)
		return err
	}

	log.Printf("[ReportByIDs] Retrieved %d products from DB", len(products))

	log.Println("[ReportByIDs] Generating Excel...")
	if err := ExportProductsToExcel(products); err != nil {
		log.Printf("[ReportByIDs] Excel generation failed: %v", err)
		return err
	}

	log.Printf("[ReportByIDs] Report successfully generated in %v", time.Since(start))
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

	const numWorkers = 5
	createDelay := 250 * time.Millisecond
	fetchDelay := 250 * time.Millisecond

	var (
		wg             sync.WaitGroup
		taskChan       = make(chan task)
		resultChan     = make(chan model.ProductPinItem)
		retryTasks     = make([]task, 0)
		retryFetchOnly = make([]task, 0)
		mu             sync.Mutex
	)

	limiter := rate.NewLimiter(20, 1)

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

				time.Sleep(createDelay) // before CreateTX

				if err := limiter.Wait(ctx); err != nil {
					log.Printf("[ERROR] Rate limiter: %v", err)
					continue
				}

				err := retryOn429(ctx, func() error {
					return utils.CreateDTOneTransaction(ctx, t.ExternalID, t.LineItem.ProductID, req.MobileNumber)
				})
				if err != nil {
					log.Printf("[WARN] CreateTX failed after retries: %v", err)
					mu.Lock()
					retryTasks = append(retryTasks, t)
					mu.Unlock()
					continue
				}

				// if err := utils.CreateDTOneTransaction(ctx, t.ExternalID, t.LineItem.ProductID, req.MobileNumber); err != nil {
				// 	log.Printf("[WARN] [Worker %d] CreateTX failed for %s: %v", workerID, t.ExternalID, err)
				// 	mu.Lock()
				// 	retryTasks = append(retryTasks, t)
				// 	mu.Unlock()
				// 	continue
				// }

				time.Sleep(fetchDelay) // before FetchTX

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
		time.Sleep(createDelay) // before CreateTX

		if err := utils.CreateDTOneTransaction(ctx, t.ExternalID, t.LineItem.ProductID, req.MobileNumber); err != nil {
			log.Printf("[ERROR] Final CreateTX failed for %s: %v", t.ExternalID, err)
			continue
		}
		retryFetchOnly = append(retryFetchOnly, t)
	}

	for _, t := range retryFetchOnly {

		log.Printf("[RETRY] Fetching transaction post-retry for ExternalID: %s", t.ExternalID)

		time.Sleep(fetchDelay) // before FetchTX

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

func retryOn429(ctx context.Context, fn func() error) error {
	backoff := time.Second
	for i := 0; i < 5; i++ {
		err := fn()
		if err == nil {
			return nil
		}
		if strings.Contains(err.Error(), "429") {
			log.Printf("[RETRY] 429 received. Backing off for %v...", backoff)
			time.Sleep(backoff)
			backoff *= 2 // exponential backoff
			continue
		}
		return err
	}
	return fmt.Errorf("failed after retries")
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
