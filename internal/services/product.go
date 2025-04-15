package services

import (
	"context"
	"log"

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
