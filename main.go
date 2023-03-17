package main

import (
	"jwt-project/common/env"
	"jwt-project/controller"
	"jwt-project/repository"
	"jwt-project/service"

	"jwt-project/routes"
)

func main() {
	env.Load()
	router := setupAllDependencies()

	router.Run(env.URL)
}

func setupAllDependencies() routes.IRouter {
	repository := repository.NewRepository()
	service := service.NewService(repository)
	controller := controller.NewController(service)
	router := routes.NewRouter(controller)

	return router
}
