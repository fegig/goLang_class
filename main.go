package main

import (
	user "web-service/app/controllers"
	middleware "web-service/app/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	router.POST("/register", user.RegisterUser)
	router.POST("/login", user.LoginUser)

	router.Use(middleware.AuthMiddleware())
	router.GET("/users", user.GetUsers)
	router.GET("/user/:userID", user.GetUser)
	router.PUT("/update/:userID", user.UpdateUser)
	router.Run("localhost:8080")
}
