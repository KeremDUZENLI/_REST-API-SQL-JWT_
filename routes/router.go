package routes

import (
	"jwt-project/controller"
	"jwt-project/middleware"

	"github.com/gin-gonic/gin"
)

type structRoutes struct {
	engine     *gin.Engine
	controller controller.InterfaceController
}

type InterfaceRoutes interface {
	Setup()
	Run(string)
}

func NewRouter(cIC controller.InterfaceController) InterfaceRoutes {
	return &structRoutes{controller: cIC}
}

func (r *structRoutes) Run(serverHost string) {
	r.engine.Run(serverHost)
}

/*
func Setup() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())

	AuthenticationRoutes(router)
	PersonRoutes(router)

	return router
}
*/

func (r *structRoutes) Setup() {
	r.engine = gin.New()

	r.engine.Use(gin.Logger())

	r.authenticationRoutes()
	r.personRoutes()
}

/*
func AuthenticationRoutes(routes *gin.Engine) {
	routes.POST("/person/signup", controller.SignUp)
	routes.POST("/person/login", controller.Login)
}

func PersonRoutes(routes *gin.Engine) {
	routes.Use(middleware.Autheticate())
	routes.GET("/person/:userid", controller.GetUser)
	routes.GET("/personall", controller.GetUsers)
}
*/

func (r *structRoutes) authenticationRoutes() {
	r.engine.POST("/person/signup", r.controller.SignUp)
	r.engine.POST("/person/login", r.controller.Login)
}

func (r *structRoutes) personRoutes() {
	r.engine.Use(middleware.Autheticate())
	r.engine.GET("/person/:userid", r.controller.GetUser)
	r.engine.GET("/personall", r.controller.GetUsers)
}
