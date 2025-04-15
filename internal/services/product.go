// // package services

// // import (
// // 	"context"
// // 	"log"
// // 	"sync"

// // 	"github.com/aakritigkmit/payment-gateway/internal/model"
// // 	"github.com/aakritigkmit/payment-gateway/internal/repository"
// // 	"github.com/aakritigkmit/payment-gateway/internal/utils"
// // 	"go.mongodb.org/mongo-driver/mongo"
// // )

// // const (
// // 	concurrencyLevel = 5
// // 	batchSize        = 100
// // )

// // func SyncProductsService(ctx context.Context, db *mongo.Database) error {
// // 	var (
// // 		offset    = 0
// // 		wg        sync.WaitGroup
// // 		productCh = make(chan []model.Product, concurrencyLevel)
// // 		errCh     = make(chan error, concurrencyLevel)
// // 		doneCh    = make(chan bool)
// // 	)

// // 	// Worker goroutines
// // 	for i := 0; i < concurrencyLevel; i++ {
// // 		wg.Add(1)
// // 		go func() {
// // 			defer wg.Done()
// // 			for products := range productCh {
// // 				if err := repository.SaveProductsToDB(ctx, db, products); err != nil {
// // 					errCh <- err
// // 				}
// // 			}
// // 		}()
// // 	}

// // 	// Chunked Fetch Loop
// // 	go func() {
// // 		for {
// // 			log.Printf("Fetching products batch offset=%d", offset)
// // 			products, err := utils.FetchProductsChunk(offset, batchSize)
// // 			if err != nil {
// // 				errCh <- err
// // 				break
// // 			}
// // 			if len(products) == 0 {
// // 				break
// // 			}
// // 			productCh <- products
// // 			offset += batchSize
// // 		}
// // 		close(productCh)
// // 	}()

// // 	// Wait + collect errors
// // 	go func() {
// // 		wg.Wait()
// // 		close(errCh)
// // 		doneCh <- true
// // 	}()

// // 	select {
// // 	case <-doneCh:
// // 		log.Println("✅ All products synced successfully")
// // 		return nil
// // 	case err := <-errCh:
// // 		return err
// // 	}
// // }

// package services

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"sync"

// 	"github.com/aakritigkmit/payment-gateway/internal/model"
// 	"github.com/aakritigkmit/payment-gateway/internal/repository"
// 	"github.com/aakritigkmit/payment-gateway/internal/utils"
// )

// type ProductService struct {
// 	repo *repository.ProductRepo
// }

// func NewProductService(repo *repository.ProductRepo) *ProductService {
// 	return &ProductService{repo}
// }

// func (s *ProductService) SyncProductsFromDTOne(ctx context.Context) error {
// 	username := os.Getenv("DTONE_USERNAME")
// 	password := os.Getenv("DTONE_PASSWORD")

// 	products, err := utils.FetchDTOneProducts(ctx, username, password)
// 	if err != nil {
// 		return fmt.Errorf("failed to fetch products from DT One: %w", err)
// 	}

// 	// Channel-based concurrency
// 	const workerCount = 5
// 	productChan := make(chan model.Product, len(products))
// 	var wg sync.WaitGroup

// 	// Start workers
// 	for i := 0; i < workerCount; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			var batch []model.Product

// 			for product := range productChan {
// 				batch = append(batch, product)
// 				if len(batch) >= 50 {
// 					_ = s.repository.SaveProductsToDB(ctx, batch) // Handle/log errors if needed
// 					batch = nil
// 				}
// 			}

// 			if len(batch) > 0 {
// 				_ = s.repo.SaveProductsToDB(ctx, batch)
// 			}
// 		}()
// 	}

// 	// Feed products into channel
// 	for _, product := range products {
// 		productChan <- product
// 	}
// 	close(productChan)

// 	wg.Wait()
// 	return nil
// }

package services

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/aakritigkmit/payment-gateway/internal/model"
	"github.com/aakritigkmit/payment-gateway/internal/repository"
	"github.com/aakritigkmit/payment-gateway/internal/utils"
)

const fetchChunkSize = 100

type ProductService struct {
	productRepo *repository.ProductRepo
}

func NewProductService(productRepo *repository.ProductRepo) *ProductService {
	return &ProductService{productRepo}
}

func (s *ProductService) SyncProducts(ctx context.Context) error {
	username := os.Getenv("DTONE_USERNAME")
	password := os.Getenv("DTONE_PASSWORD")

	allProducts, err := utils.FetchProductsFromDTOneAPI(username, password)
	if err != nil {
		return err
	}

	total := len(allProducts)
	log.Printf("Total products to sync: %d", total)

	ch := make(chan []model.Product)
	wg := sync.WaitGroup{}

	// Workers to consume from channel
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range ch {
				for _, product := range batch {
					if err := s.productRepo.SaveOrUpdateProduct(ctx, product); err != nil {
						log.Printf("Error saving product_id=%s: %v", product.ProductId, err)
					}
				}
			}
		}()
	}

	// Feed channel
	for i := 0; i < total; i += fetchChunkSize {
		end := i + fetchChunkSize
		if end > total {
			end = total
		}
		ch <- allProducts[i:end]
	}
	close(ch)

	wg.Wait()
	log.Println("✅ DT One product sync completed.")
	return nil
}
