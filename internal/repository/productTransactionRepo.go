package repository

import (
	"context"
	"fmt"

	"github.com/aakritigkmit/payment-gateway/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductTransactionRepo struct {
	collection *mongo.Collection
}

func NewProductTransactionRepo(db *mongo.Database) *ProductTransactionRepo {
	return &ProductTransactionRepo{
		collection: db.Collection("product_transactions_ak"),
	}
} 

// SaveProductTransaction inserts a new product transaction into the database
func (r *ProductTransactionRepo) SaveProductTransaction(ctx context.Context, tx model.ProductTransaction) error {
	_, err := r.collection.InsertOne(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to save product transaction: %w", err)
	}
	return nil
}

// FindProductTransactionByExternalID finds a product transaction by its external_id
func (r *ProductTransactionRepo) FindProductTransactionByExternalID(ctx context.Context, externalID string) (*model.ProductTransaction, error) {
	var tx model.ProductTransaction
	filter := bson.M{"external_id": externalID}

	err := r.collection.FindOne(ctx, filter).Decode(&tx)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to find product transaction by external_id: %w", err)
	}

	return &tx, nil
}
