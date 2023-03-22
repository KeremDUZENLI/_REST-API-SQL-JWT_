package controller

import (
	"jwt-project/common/env"
	"jwt-project/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
