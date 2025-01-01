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

	// start the server
	router.Run(":" + port)
}