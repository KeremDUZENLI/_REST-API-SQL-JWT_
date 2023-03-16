package routes

import (
	"jwt-project/controller"
	"jwt-project/middleware"

	"github.com/gin-gonic/gin"
)

type structRouter struct {
	engine     *gin.Engine
	controller controller.IController
}

type IRouter interface {
	Run(string)
}

func NewRouter(cIC controller.IController) IRouter {
	router := &structRouter{controller: cIC}
	router.setup()

	return router
}

func (r *structRouter) Run(serverHost string) {
	r.engine.Run(serverHost)
}

func (r *structRouter) setup() {
	r.engine = gin.New()

	r.engine.Use(gin.Logger())

	r.authenticationRoutes()
	r.personRoutes()
}

func (r *structRouter) authenticationRoutes() {
	r.engine.POST("/person/signup", r.controller.SignUp)
	r.engine.POST("/person/login", r.controller.Login)
}

func (r *structRouter) personRoutes() {
	r.engine.Use(middleware.Autheticate())
	r.engine.GET("/person/:userid", r.controller.GetUser)
	r.engine.GET("/personall", r.controller.GetUsers)
}
