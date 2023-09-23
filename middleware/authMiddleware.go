package middleware

import (
	"jwt-project/database/model"
	"jwt-project/middleware/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if tokenIsExist(c, clientToken) {
			return
		}

		claims, msg := token.ValidateToken(clientToken)
		if tokenAuthenticate(c, msg) {
			return
		}

		setContextClaims(c, claims)
		c.Next()
	}
}

func tokenIsExist(c *gin.Context, token string) bool {
	if token == model.NONE {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization header provided"})
		c.Abort()
		return true
	}

	return false
}

func tokenAuthenticate(c *gin.Context, message string) bool {
	if message != model.NONE {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		c.Abort()
		return true
	}

	return false
}

func setContextClaims(c *gin.Context, claims *token.SignedDetails) {
	c.Set("first_name", claims.FirstName)
	c.Set("last_name", claims.LastName)
	c.Set("email", claims.Email)
	c.Set("usertype", claims.UserType)
	c.Set("uid", claims.Uid)
	c.Next()
}
