package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aakritigkmit/payment-gateway/internal/dto"
	"github.com/aakritigkmit/payment-gateway/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepo struct {
	collection *mongo.Collection
}

func NewOrderRepo(db *mongo.Database) *OrderRepo {
	return &OrderRepo{
		collection: db.Collection("orders"),
	}
}

func (r *OrderRepo) SaveOrder(ctx context.Context, order model.Order) error {
	_, err := r.collection.InsertOne(ctx, order)
	return err
}

func (r *OrderRepo) UpdateOrder(referenceID string, payload *dto.UpdateOrderPayload) error {
	if referenceID == "" {
		return fmt.Errorf("empty transaction reference ID")
	}

	update := bson.M{}
	if payload.Status != "" {
		update["status"] = payload.Status
	}
	if payload.Amount != 0 {
		update["amount"] = payload.Amount
	}
	if payload.Currency != "" {
		update["currency"] = payload.Currency
	}
	if payload.UserID != "" {
		update["userId"] = payload.UserID
	}

	update["updatedAt"] = time.Now()

	result, err := r.collection.UpdateOne(
		context.TODO(),
		bson.M{"transactionReferenceId": referenceID},
		bson.M{"$set": update},
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("order with transactionReferenceId %s not found", referenceID)
	}

	return nil
}

func (r *OrderRepo) GetOrderByID(ctx context.Context, orderID string) (model.Order, error) {
	var order model.Order

	err := r.collection.FindOne(ctx, bson.M{"transactionReferenceId": orderID}).Decode(&order)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
