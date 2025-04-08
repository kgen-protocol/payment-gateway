package repository

import (
	"context"

	"github.com/aakritigkmit/payment-gateway/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	Client *mongo.Collection
}

func NewUserRepo(client *mongo.Collection) *UserRepo {
	return &UserRepo{
		Client: client,
	}
}

func (c *UserRepo) CreateUser(user model.User) (string, error) {
	result, err := c.Client.InsertOne(context.Background(), user)
	if err != nil {
		return "", nil
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
	//
	// Hex returns the hex encoding of the ObjectID as a string.
}
