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
	//tmpl読み込み
	r.LoadHTMLGlob("view/*tmpl")

	port := os.Getenv("PORT")
	http.HandleFunc("/", hello)
	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "index.tmpl")
	})
	http.ListenAndServe(":"+port, nil)
}
