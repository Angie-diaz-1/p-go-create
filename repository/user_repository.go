package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"usersProject/models"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error)
}

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database, collectionName string) *UserRepository {
	return &UserRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error) {
	return r.collection.InsertOne(ctx, user)
}
