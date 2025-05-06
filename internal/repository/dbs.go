package repository

import (
	"context"

	"github.com/aakritigkmit/payment-gateway/internal/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type DBSRepo struct {
	collection *mongo.Collection
}

func NewDBSRepo(db *mongo.Database) *DBSRepo {
	return &DBSRepo{
		collection: db.Collection("dbs_bank_statements"),
	}
}
func (r *DBSRepo) SaveBankStatement(ctx context.Context, resp *model.Camt053Response) error {
	_, err := r.collection.InsertOne(ctx, resp)
	return err
}
