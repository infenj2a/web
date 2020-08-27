// main.go

package main

import (
	"github.com/gin-gonic/gin"
	"main/controller"
	"os"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("view/*.html")

	r.GET("/", controller.HelloPage)
	r.GET("/get", controller.PostPage)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
