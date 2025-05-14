package service

import (
	"context"
	"testing"
	"usersProject/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mock minimal que implementa SOLO Create del repositorio
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

// Test de CreateUser en el servicio
func TestCreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	user := &models.User{
		Name:     "Test",
		LastName: "User",
		Email:    "test@example.com",
		Password: "1234",
	}

	expected := &mongo.InsertOneResult{InsertedID: "fakeID123"}
	mockRepo.On("Create", mock.Anything, user).Return(expected, nil)

	result, err := userService.CreateUser(context.TODO(), user)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}
