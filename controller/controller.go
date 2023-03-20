package controller

import (
	"jwt-project/common/env"
	"jwt-project/dto"
	"jwt-project/service"
	"net/http"

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

func (sC user) SignUp(c *gin.Context) {
	var dtoPerson dto.DtoSignUp
	if err := c.BindJSON(&dtoPerson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't bind!"})
		return
	}

	insert, err := sC.service.CreateUser(c, dtoPerson)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, insert)
}

func (sC user) LogIn(c *gin.Context) {
	var dtoPerson dto.DtoLogIn
	if err := c.BindJSON(&dtoPerson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't bind!"})
		return
	}

	foundPerson, err := sC.service.FindUser(c, dtoPerson)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &foundPerson.ID)
}

func (sC user) GetUser(c *gin.Context) {
	var dtoPerson dto.GetUser

	personId := c.Param("userId")

	person, err := sC.service.GetUserByID(c, dtoPerson, personId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, person)
}

func (sC user) GetUsers(c *gin.Context) {
	var allUsers env.LISTUSERS

	allUsersList, err := sC.service.GetAllUsers(c, allUsers)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, allUsersList)
}
