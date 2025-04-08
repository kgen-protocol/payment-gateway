package repository

import (
	"context"

	"github.com/aakritigkmit/payment-gateway/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepo struct {
	Client *mongo.Collection
}

func NewAuthRepo(client *mongo.Collection) *AuthRepo {
	return &AuthRepo{
		Client: client,
	}
}

func (r *AuthRepo) FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.Client.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
