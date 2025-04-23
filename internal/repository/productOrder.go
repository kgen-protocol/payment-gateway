package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aakritigkmit/payment-gateway/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductOrderRepo struct {
	collection               *mongo.Collection
	productPinDumpcollection *mongo.Collection
}

func NewProductOrderRepo(db *mongo.Database) *ProductOrderRepo {
	return &ProductOrderRepo{
		collection:               db.Collection("product_order_ab"),
		productPinDumpcollection: db.Collection("product_pin_dump"),
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

func (r *ProductOrderRepo) UpdateProductPinsByOrderID(ctx context.Context, orderId string, update bson.M) error {
	filter := bson.M{"orderID": orderId}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *ProductOrderRepo) SaveProductPinsDump(ctx context.Context, dump model.ProductPinDump) error {
	_, err := r.productPinDumpcollection.InsertOne(ctx, dump)
	return err
}

func (r *ProductOrderRepo) GetPinsByOrderID(ctx context.Context, orderId string) ([]model.ProductPinItem, error) {
	var result model.ProductPinDump
	err := r.productPinDumpcollection.FindOne(ctx, bson.M{"orderID": orderId}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.ProductPins, nil
}

func (r *ProductOrderRepo) UpdateProductOrderWithPins(ctx context.Context, orderId string, pins []model.ProductPinItem) error {
	update := bson.M{
		"$set": bson.M{
			"productPins": pins,
			"updated_at":  time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"orderID": orderId}, update)
	return err
}

func (r *ProductOrderRepo) UpdateProductPins(ctx context.Context, orderID string, pins []model.ProductPinItem) error {

	filter := bson.M{"orderID": orderID}
	update := bson.M{
		"$set": bson.M{
			"productPins": pins,
			"updated_at":  time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update product pins for orderID %s: %w", orderID, err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found for orderID %s", orderID)
	}

	return nil
}
