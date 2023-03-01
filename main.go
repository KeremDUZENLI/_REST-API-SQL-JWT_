package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"jwt-project/routes"
)

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthenticationRoutes(router)
	routes.PersonRoutes(router)

	router.Run(":" + port)
}
