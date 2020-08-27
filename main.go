// main.go

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"main/controller"
	"os"
)

func main() {
	db, _ := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	defer db.Close()

	r := gin.Default()
	r.LoadHTMLGlob("view/*.html")

	r.GET("/", controller.HelloPage)
	r.POST("/", controller.PostPage)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
