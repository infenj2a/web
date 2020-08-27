// main.go

package main

import (
	"github.com/gin-gonic/gin"
	"main/controller"
	"main/util"
	"os"
)

func main() {
	server := controller.Server{
		DB: util.InitDB(),
	}

	r := gin.Default()
	r.LoadHTMLGlob("view/*.html")

	r.GET("/", server.HelloPage)
	r.POST("/", server.PostPage)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
