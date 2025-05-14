package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"usersProject/models"
	"usersProject/service"
)

type UserController struct {
	service service.UserServiceInterface
}

func NewUserController(s service.UserServiceInterface) *UserController {
	return &UserController{service: s}
}

func (uc *UserController) RegisterRoutes(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.POST("/", uc.CreateUser)

	}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := uc.service.CreateUser(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}
