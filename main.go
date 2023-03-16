package main

import (
	"os"

	"jwt-project/common/env"
	"jwt-project/controller"
	"jwt-project/repository"
	"jwt-project/service"

	"jwt-project/routes"
)

func main() {
	env.Load()
	router := setupAllDependencies()

	port := os.Getenv("PORT")
	url := ":" + port

	router.Run(url)
}

func setupAllDependencies() routes.IRouter {
	repository := repository.NewRepository()
	service := service.NewService(repository)
	controller := controller.NewController(service)
	router := routes.NewRouter(controller)

	return router
}
