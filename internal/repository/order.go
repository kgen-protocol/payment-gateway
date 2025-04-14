package repository

import (
	"context"
	"fmt"
	"time"

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

func (r *OrderRepo) SaveOrder(ctx context.Context, order *model.Order) error {
	_, err := r.collection.InsertOne(ctx, order)
	return err
}

// func (r *OrderRepo) UpdateOrder(referenceID string, payload *dto.UpdateOrderPayload) error {
// 	if referenceID == "" {
// 		return fmt.Errorf("empty transaction reference ID")
// 	}

// 	update := bson.M{}
// 	if payload.Status != "" {
// 		update["status"] = payload.Status
// 	}
// 	if payload.Amount != 0 {
// 		update["amount"] = payload.Amount
// 	}
// 	if payload.Currency != "" {
// 		update["currency"] = payload.Currency
// 	}
// 	if payload.UserID != "" {
// 		update["userId"] = payload.UserID
// 	}

// 	update["updatedAt"] = time.Now()

// 	result, err := r.collection.UpdateOne(
// 		context.TODO(),
// 		bson.M{"transactionReferenceId": referenceID},
// 		bson.M{"$set": update},
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	if result.MatchedCount == 0 {
// 		return fmt.Errorf("order with transactionReferenceId %s not found", referenceID)
// 	}

// 	return nil
// }

func (r *OrderRepo) UpdateOrder(ctx context.Context, orderID string, updatedOrder *model.Order) error {
	// Update the UpdatedAt field
	updatedOrder.UpdatedAt = time.Now()

	// Convert the updatedOrder struct to a bson.M map
	updateData, err := toBsonM(updatedOrder)
	if err != nil {
		return fmt.Errorf("failed to convert order to bson: %w", err)
	}

	// Prepare the update query
	update := bson.M{
		"$set": updateData,
	}

	// Run the update
	filter := bson.M{"order_id": orderID}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}

func toBsonM(doc interface{}) (bson.M, error) {
	data, err := bson.Marshal(doc)
	if err != nil {
		return nil, err
	}
	var result bson.M
	err = bson.Unmarshal(data, &result)
	return result, err
}
