package controller

import (
	"jwt-project/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
