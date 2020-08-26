// main.go

package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World!")
}

func main() {
	router := gin.Default()
	port := os.Getenv("PORT")
	http.HandleFunc("/", hello)
	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello!!")
	})
	http.ListenAndServe(":"+port, nil)
}
