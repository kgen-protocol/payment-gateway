package repository

import (
	"context"
	"log"
	"reflect"

	"github.com/aakritigkmit/payment-gateway/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepo struct {
	collection *mongo.Collection
}

func NewProductRepo(db *mongo.Database) *ProductRepo {
	return &ProductRepo{collection: db.Collection("products")}
}

func (r *ProductRepo) SaveOrUpdateProduct(ctx context.Context, product model.Product) error {
	var existing model.Product
	err := r.collection.FindOne(ctx, bson.M{"product_id": product.ProductId}).Decode(&existing)

	if err == mongo.ErrNoDocuments {
		product.ID = primitive.NewObjectID()
		_, err := r.collection.InsertOne(ctx, product)
		if err != nil {
			log.Printf("❌ Insert failed for product_id=%s: %v", product.ProductId, err)
		}
		return err
	} else if err != nil {
		return err
	}

	if reflect.DeepEqual(existing, product) {
		log.Printf("⏩ Skipped identical product_id=%s", product.ProductId)
		return nil
	}

	_, err = r.collection.ReplaceOne(ctx, bson.M{"product_id": product.ProductId}, product)
	if err != nil {
		log.Printf("❌ Update failed for product_id=%s: %v", product.ProductId, err)
	}
	return err
}
