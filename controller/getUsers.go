package controller

import (
	"jwt-project/common/env"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (sC user) GetUsers(c *gin.Context) {
	var allUsers env.LISTUSERS

	allUsersList, err := sC.service.GetAllUsers(c, allUsers)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, allUsersList)
}
