package main

import (
	"jwt/controllers"
	"jwt/initilizers"
	"jwt/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	initilizers.Loadenv()

	initilizers.ConnectToDB()

	initilizers.SyncDB()

	r := gin.Default()

	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	r.Run(":3000")

}
