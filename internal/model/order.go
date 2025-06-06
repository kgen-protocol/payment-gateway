package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID                 string             `bson:"userId"`
	TransactionReferenceId string             `bson:"transactionReferenceId"`
	Amount                 float32            `bson:"amount"`
	Currency               string             `bson:"currency"`
	Status                 string             `bson:"status"`
	CreatedAt              time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt              time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
