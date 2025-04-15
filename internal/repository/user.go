package repository

import (
	"context"
	"fmt"

	"github.com/aakritigkmit/payment-gateway/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	db *mongo.Collection
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{db: db.Collection("users")}
}

func (c *UserRepo) CreateUser(user model.User) (string, error) {
	result, err := c.db.InsertOne(context.Background(), user)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *UserRepo) FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, err
	}
	return &user, nil
}
