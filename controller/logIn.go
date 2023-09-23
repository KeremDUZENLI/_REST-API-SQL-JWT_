package controller

import (
	"jwt-project/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
