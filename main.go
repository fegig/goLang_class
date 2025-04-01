package main

import (
	user "github.com/fegig/goLang_class/app/controllers"
	middleware "github.com/fegig/goLang_class/app/middlewares"

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
