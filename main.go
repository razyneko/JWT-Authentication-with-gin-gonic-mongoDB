package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/razyneko/jwt-auth-with-go-gin-gonic-mongodb/routes"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error laoding .env file")
	}
	port := os.Getenv("PORT")

	if port == ""{
		port="8000"
	}

	//initializing router
	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	// APIs
	router.GET("/api-1", func(c *gin.Context){
		// *gin.Context has w http.ResposneWriter and r *http.Request built in
		c.JSON(200, gin.H{"success":"Access granted for api-1"})
	})

	router.GET("/api-2", func(c *gin.Context){
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.Run(":" + port)
}