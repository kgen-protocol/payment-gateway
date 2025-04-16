package repository

import (
	"context"
	"time"

	"github.com/aakritigkmit/payment-gateway/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type RefundRepo struct {
	db *mongo.Database
}

func NewRefundRepository(db *mongo.Database) *RefundRepo {
	return &RefundRepo{db: db}
}

func (r *RefundRepo) SaveRefund(ctx context.Context, refund model.Refund) error {
	collection := r.db.Collection(" order_refunds")
	refund.CreatedAt = time.Now()

	_, err := collection.InsertOne(ctx, refund)
	return err
}
