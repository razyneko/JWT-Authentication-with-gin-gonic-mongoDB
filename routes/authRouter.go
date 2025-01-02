package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/razyneko/jwt-auth-with-go-gin-gonic-mongodb/controllers"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	// incomingRoutes -> This is the main router instance where the routes are registered.

	// these are not protected routes since user doesnt have a token before signup
	
	incomingRoutes.POST("users/signup", controllers.SignUp())
	// Handles user signup requests. The client sends data and the SignUp function processes this data, storing it in the database after validation.

	incomingRoutes.POST("users/login", controllers.Login())
	// Handles user login requests. The client sends credentials and the Login function verifies the details and generates a JWT if valid.
}
