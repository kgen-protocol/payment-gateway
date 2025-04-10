package repository

import (
	"context"

	"github.com/aakritigkmit/payment-gateway/internal/model"
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
