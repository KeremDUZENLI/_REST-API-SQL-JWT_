package routes

import (
	"jwt-project/controller"
	"jwt-project/middleware"

	"github.com/gin-gonic/gin"
)

type router struct {
	engine     *gin.Engine
	controller controller.User
}

type Router interface {
	Run(string)
}

func NewRouter(cIC controller.User) Router {
	router := &router{controller: cIC}
	router.setup()

	return router
}

func (r *router) Run(serverHost string) {
	r.engine.Run(serverHost)
}

func (r *router) setup() {
	r.engine = gin.New()

	r.engine.Use(gin.Logger())

	r.authenticationRoutes()
	r.personRoutes()
}

func (r *router) authenticationRoutes() {
	r.engine.POST("/person/signup", r.controller.SignUp)
	r.engine.POST("/person/login", r.controller.LogIn)
}

func (r *router) personRoutes() {
	r.engine.Use(middleware.Authenticate())
	r.engine.GET("/person/:userId", r.controller.GetUser)
	r.engine.GET("/personall", r.controller.GetUsers)
}
