package repository

import (
	"context"
	"fmt"

	"github.com/aakritigkmit/payment-gateway/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductOrderRepo struct {
	collection *mongo.Collection
}

func NewProductOrderRepo(db *mongo.Database) *ProductOrderRepo {
	return &ProductOrderRepo{
		collection: db.Collection("product_order"),
	}
}

// SaveProductPins saves a new order with product pins to the database
func (r *ProductOrderRepo) SaveProductPins(ctx context.Context, pin model.ProductPin) error {
	_, err := r.collection.InsertOne(ctx, pin)
	if err != nil {
		return fmt.Errorf("failed to save product pins: %w", err)
	}
	return nil
}
