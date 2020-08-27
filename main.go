// main.go

package main

import (
	"github.com/gin-gonic/gin"
	"main/controller"
	"os"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("view/*.html")

	router.GET("/", controller.HelloPage)
	router.GET("/get", controller.PostPage)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
