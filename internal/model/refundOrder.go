package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Refund struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty"`
	RefundID               primitive.ObjectID `bson:"refund_id"`
	ParentOrderID          string             `bson:"pine_order_id"`
	MerchantOrderReference string             `bson:"merchant_order_reference"`
	RefundAmount           interface{}        `bson:"refund_amount"`
	Status                 string             `bson:"status"` // "success", "failed", etc.
	CreatedAt              time.Time          `bson:"created_at"`
}
