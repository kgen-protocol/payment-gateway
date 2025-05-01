package repository

import (
	"context"
	"fmt"

	"github.com/aakritigkmit/payment-gateway/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepo struct {
	collection *mongo.Collection
}

func NewProductRepo(db *mongo.Database) *ProductRepo {
	return &ProductRepo{collection: db.Collection("products")}
}

func (r *ProductRepo) FindOrCreateProduct(ctx context.Context, product model.Product) error {
	filter := bson.M{"unique_id": product.UniqueId}

	// Check if the product already exists
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to check for existing product: %w", err)
	}

	// If it exists, skip insert
	if count > 0 {
		return nil
	}

	// If not, insert the new product
	_, err = r.collection.InsertOne(ctx, product)
	if err != nil {
		return fmt.Errorf("failed to insert product: %w", err)
	}

	return nil
}

func (r *ProductRepo) GetProductsByUniqueIDs(ctx context.Context, ids []int) ([]model.Product, error) {
	filter := bson.M{"unique_id": bson.M{"$in": ids}}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []model.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}
