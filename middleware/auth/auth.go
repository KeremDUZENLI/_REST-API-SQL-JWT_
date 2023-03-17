package auth

import (
	"errors"

	"github.com/gin-gonic/gin"

	"jwt-project/database/model"
)

func CheckPersonType(c *gin.Context, role string) (err error) {
	userType := c.GetString("usertype")
	err = nil

	if userType != role {
		err = errors.New("unauthorized to access this resource")
		return err
	}

	return err
}

func MatchPersonTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("usertype")
	uid := c.GetString("uid")
	err = nil

	if userType == model.USER && uid != userId {
		err = errors.New("unauthorized to access this resource")
		return err
	}

	err = CheckPersonType(c, userType)
	return err
}
