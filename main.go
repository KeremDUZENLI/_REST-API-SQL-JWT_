package main

import (
	"os"

	"github.com/joho/godotenv"

	"jwt-project/routes"
)

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")

	routes.Setup().Run(":" + port)
}
