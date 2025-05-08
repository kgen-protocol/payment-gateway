package repository

import (
	"context"

	"github.com/aakritigkmit/payment-gateway/internal/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type DBSRepo struct {
	bankStatementcollection            *mongo.Collection
	bankIntradayNotificationCollection *mongo.Collection
	bankIncomingNotificationCollection *mongo.Collection
}

func NewDBSRepo(db *mongo.Database) *DBSRepo {
	return &DBSRepo{
		bankStatementcollection:            db.Collection("dbs_bank_statements"),
		bankIntradayNotificationCollection: db.Collection("dbs_intraday_bank_notifications"),
		bankIncomingNotificationCollection: db.Collection("dbs_incoming_bank_notifications"),
	}
}

func (r *DBSRepo) SaveBankStatement(ctx context.Context, req model.CAMT053Request) error {
	_, err := r.bankStatementcollection.InsertOne(ctx, req)
	return err
}

func (r *DBSRepo) SaveIntradayNotification(ctx context.Context, payload model.IntradayNotificationPayload) error {
	_, err := r.bankIntradayNotificationCollection.InsertOne(ctx, payload)
	return err
}

func (r *DBSRepo) SaveIncomingNotification(ctx context.Context, payload model.IncomingNotificationPayload) error {
	_, err := r.bankIncomingNotificationCollection.InsertOne(ctx, payload)
	return err
}
