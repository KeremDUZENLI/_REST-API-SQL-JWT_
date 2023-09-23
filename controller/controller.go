package controller

import (
	"jwt-project/service"

	"github.com/gin-gonic/gin"
)

type user struct {
	service service.MongoService
}

type User interface {
	SignUp(c *gin.Context)
	LogIn(c *gin.Context)
	GetUser(c *gin.Context)
	GetUsers(c *gin.Context)
}

func NewUser(sIS service.MongoService) User {
	return &user{sIS}
}
