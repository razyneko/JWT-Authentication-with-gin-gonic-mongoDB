package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/razyneko/jwt-auth-with-go-gin-gonic-mongodb/controllers"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	// these are not protected routes since user doesnt have a token in this
	
	incomingRoutes.POST("users/signup", controllers.SignUp())
	incomingRoutes.POST("users/login", controllers.Login())
}
