package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID                 string             `bson:"user_id"`
	MerchantOrderReference int                `bson:"merchant_order_reference"`
	Amount                 float32            `bson:"amount"`
	Currency               float32            `bson:"currency"`
	Status                 string             `bson:"status"` // pending, success, failure
	CreatedAt              time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt              time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
