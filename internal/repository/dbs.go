package repository

import (
	"context"

	"github.com/aakritigkmit/payment-gateway/internal/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type DBSRepo struct {
	bankStatementcollection    *mongo.Collection
	bankNotificationCollection *mongo.Collection
}

func NewDBSRepo(db *mongo.Database) *DBSRepo {
	return &DBSRepo{
		bankStatementcollection:    db.Collection("dbs_bank_statements"),
		bankNotificationCollection: db.Collection("dbs_bank_notifications"),
	}
}

func (r *DBSRepo) SaveBankStatement(ctx context.Context, req model.CAMT053Request) error {
	_, err := r.bankStatementcollection.InsertOne(ctx, req)
	return err
}

func (r *DBSRepo) SaveNotification(ctx context.Context, req model.NotificationPayload) error {
	_, err := r.bankNotificationCollection.InsertOne(ctx, req)
	return err
}
