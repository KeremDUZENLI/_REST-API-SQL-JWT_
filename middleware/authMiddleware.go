package middleware

import (
	"jwt-project/middleware/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Autheticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")

		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization header provided"})
			c.Abort()
			return
		}

		claims, _ := token.ValidateToken(clientToken)

		setContextClaims(c, claims)
		c.Next()
	}
}

func setContextClaims(c *gin.Context, claims *token.SignedDetails) {
	c.Set("first_name", claims.FirstName)
	c.Set("last_name", claims.LastName)
	c.Set("email", claims.Email)
	c.Set("usertype", claims.UserType)
	c.Set("uid", claims.Uid)
	c.Next()
}
