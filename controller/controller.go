package controller

import (
	"jwt-project/database/model"
	"jwt-project/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SignUp(c *gin.Context) {
	var person model.Person
	c.BindJSON(&person)

	insert, err := service.InsertInDatabase(c, person)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, insert)
}

func Login(c *gin.Context) {
	var person model.Person
	c.BindJSON(&person)

	foundPerson, err := service.FindInDatabase(c, person)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &foundPerson)
}

func GetUser(c *gin.Context) {
	var person model.Person

	personId := c.Param("userid")

	person, err := service.GetFromDatabase(c, person, personId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, person)
}

func GetUsers(c *gin.Context) {
	var allUsers []primitive.M

	allUsersList, err := service.GetallFromDatabase(c, allUsers)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, allUsersList)
}
