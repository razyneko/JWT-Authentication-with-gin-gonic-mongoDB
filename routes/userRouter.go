package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/razyneko/jwt-auth-with-go-gin-gonic-mongodb/controllers"
	"github.com/razyneko/jwt-auth-with-go-gin-gonic-mongodb/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate()) // to ensure both routes are protected
	// after login, user will have a token
	// user route can't be used without JWT
	// private routes require authentication to use
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:user_id", controllers.GetUser())
}
