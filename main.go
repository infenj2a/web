// main.go

package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("view/*.html")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{})
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	router.Run(":" + port)
}
