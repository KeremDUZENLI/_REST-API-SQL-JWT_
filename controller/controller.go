package controller

import (
	"jwt-project/service"

	"github.com/gin-gonic/gin"
)

func SignUp() gin.HandlerFunc {
	return service.InsertInDatabase
}

func Login() gin.HandlerFunc {
	return service.FindInDatabase
}

func GetUser() gin.HandlerFunc {
	return service.GetFromDatabase
}

func GetUsers() gin.HandlerFunc {
	return service.GetallFromDatabase
}
