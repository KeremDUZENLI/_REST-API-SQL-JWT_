package controller

import (
	"jwt-project/dto"
	"jwt-project/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type structController struct {
	service service.InterfaceService
}

type InterfaceController interface {
	SignUp(c *gin.Context)
	Login(c *gin.Context)
	GetUser(c *gin.Context)
	GetUsers(c *gin.Context)
}

func NewController(sIS service.InterfaceService) InterfaceController {
	return &structController{service: sIS}
}

func (sC structController) SignUp(c *gin.Context) {
	var dtoPerson dto.DtoSignUp
	c.BindJSON(&dtoPerson)

	insert, err := sC.service.InsertInDatabase(c, dtoPerson)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, insert)
}

/*
func SignUp(c *gin.Context) {
	var dtoPerson dto.DtoSignUp
	c.BindJSON(&dtoPerson)

	insert, err := service.InsertInDatabase(c, dtoPerson)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, insert)
}
*/

func (sC structController) Login(c *gin.Context) {
	var dtoPerson dto.DtoLogIn
	c.BindJSON(&dtoPerson)

	foundPerson, err := sC.service.FindInDatabase(c, dtoPerson)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &foundPerson.ID)
}

/*
func Login(c *gin.Context) {
	var dtoPerson dto.DtoLogIn
	c.BindJSON(&dtoPerson)

	foundPerson, err := service.FindInDatabase(c, dtoPerson)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &foundPerson.ID)
}
*/

func (sC structController) GetUser(c *gin.Context) {
	var dtoPerson dto.GetUser

	personId := c.Param("userid")

	person, err := sC.service.GetFromDatabase(c, dtoPerson, personId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, person)
}

/*
func GetUser(c *gin.Context) {
	var dtoPerson dto.GetUser

	personId := c.Param("userid")

	person, err := service.GetFromDatabase(c, dtoPerson, personId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, person)
}
*/

func (sC structController) GetUsers(c *gin.Context) {
	var allUsers []primitive.M

	allUsersList, err := sC.service.GetallFromDatabase(c, allUsers)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, allUsersList)
}

/*
func GetUsers(c *gin.Context) {
	var allUsers []primitive.M

	allUsersList, err := service.GetallFromDatabase(c, allUsers)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, allUsersList)
}
*/
