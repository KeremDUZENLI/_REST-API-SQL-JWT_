package routes

import (
	"github.com/gin-gonic/gin"
)

type structRoutes struct {
}

type InterfaceRoutes interface {
}

func Setup() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())

	AuthenticationRoutes(router)
	PersonRoutes(router)

	return router
}
