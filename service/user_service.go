// ✅ user_service.go (solo con CreateUser)
package service

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"usersProject/models"
	"usersProject/repository"
)

type UserServiceInterface interface {
	CreateUser(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error)
}

type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error) {
	return s.repo.Create(ctx, user)
}
