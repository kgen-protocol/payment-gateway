package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aakritigkmit/payment-gateway/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductTransactionRepo struct {
	collection *mongo.Collection
}

func NewProductTransactionRepo(db *mongo.Database) *ProductTransactionRepo {
	return &ProductTransactionRepo{
		collection: db.Collection("product_transactions_ab"),
	}
}

// SaveProductTransaction inserts a new product transaction into the database
func (r *ProductTransactionRepo) SaveProductTransaction(ctx context.Context, tx model.ProductTransaction) error {
	_, err := r.collection.InsertOne(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to save product transaction: %w", err)
	}
	return nil
}

// FindProductTransactionByExternalID finds a product transaction by its external_id
func (r *ProductTransactionRepo) FindProductTransactionByExternalID(ctx context.Context, externalID string) (*model.ProductTransaction, error) {
	var tx model.ProductTransaction
	filter := bson.M{"external_id": externalID}

	err := r.collection.FindOne(ctx, filter).Decode(&tx)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to find product transaction by external_id: %w", err)
	}

	return &tx, nil
}

func (r *ProductTransactionRepo) UpdateProductTransaction(ctx context.Context, externalID string, updatedData model.ProductTransaction) error {
	filter := bson.M{"external_id": externalID}

	// Step 1: Fetch the existing document to get created_at
	var existing model.ProductTransaction
	err := r.collection.FindOne(ctx, filter).Decode(&existing)
	if err != nil {
		return err
	}

	// Step 2: Calculate updation_time
	now := time.Now()
	updationTime := now.Sub(existing.CreatedAt) // assuming CreatedAt is of type time.Time

	// Step 3: Add updation_time to the update object
	update := bson.M{
		"$set": bson.M{
			"benefits":                     updatedData.Benefits,
			"confirmation_date":            updatedData.ConfirmationDate,
			"confirmation_expiration_date": updatedData.ConfirmationExpirationDate,
			"creation_date":                updatedData.CreationDate,
			"credit_party_identifier":      updatedData.CreditPartyIdentifier,
			"external_id":                  updatedData.ExternalID,
			"id":                           updatedData.ID,
			"operator_reference":           updatedData.OperatorReference,
			"pin":                          updatedData.Pin,
			"prices":                       updatedData.Prices,
			"product":                      updatedData.Product,
			"promotions":                   updatedData.Promotions,
			"rates":                        updatedData.Rates,
			"status":                       updatedData.Status,
			"updated_at":                   now,
			"updation_time":                updationTime.String(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}
