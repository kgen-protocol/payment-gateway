package repository

import (
	"context"

	"github.com/aakritigkmit/payment-gateway/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionRepo struct {
	collection *mongo.Collection
}

func NewTransactionRepo(db *mongo.Database) *TransactionRepo {
	return &TransactionRepo{
		collection: db.Collection("transactions"),
	}
}

func (r *TransactionRepo) SaveTransaction(ctx context.Context, transaction model.Transaction) error {
	_, err := r.collection.InsertOne(ctx, transaction)
	return err
}

func (r *TransactionRepo) UpdateTransactionByOrderID(ctx context.Context, orderID string, tx model.Transaction) error {
	filter := bson.M{"order_id": orderID}
	update := bson.M{"$set": tx}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
