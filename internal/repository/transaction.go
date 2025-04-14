package repository

import (
	"context"

	"github.com/aakritigkmit/payment-gateway/internal/model"
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

func (r *TransactionRepo) SaveTransaction(ctx context.Context, transaction *model.Transaction) error {
	_, err := r.collection.InsertOne(ctx, transaction)
	return err
}
