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
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	router.Run(":" + port)
}
