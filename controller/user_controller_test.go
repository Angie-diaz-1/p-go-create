package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"usersProject/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

// Métodos innecesarios eliminados para este microservicio

func TestCreateUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockUserService)
	userController := NewUserController(mockService)

	router := gin.Default()
	userController.RegisterRoutes(router)

	user := &models.User{
		Name:     "Juan",
		LastName: "Henao",
		Email:    "juan@example.com",
		Password: "123456",
	}

	body, _ := json.Marshal(user)
	expected := &mongo.InsertOneResult{InsertedID: "fake-id"}

	mockService.On("CreateUser", mock.Anything, user).Return(expected, nil)

	req, _ := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	mockService.AssertExpectations(t)
}

func TestCreateUserHandler_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	controller := NewUserController(mockService)
	router := gin.Default()
	controller.RegisterRoutes(router)

	invalidJSON := []byte(`{"name": "Juan", "email": "juan@example.com", "password": "123456"`)

	req, _ := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestCreateUserHandler_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	controller := NewUserController(mockService)
	router := gin.Default()
	controller.RegisterRoutes(router)

	user := &models.User{
		Name:     "Juan",
		LastName: "Henao",
		Email:    "juan@example.com",
		Password: "123456",
	}
	body, _ := json.Marshal(user)

	mockService.On("CreateUser", mock.Anything, user).Return((*mongo.InsertOneResult)(nil), errors.New("error en la creación"))

	req, _ := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	mockService.AssertExpectations(t)
}
