// main.go

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// func hello(w http.ResponseWriter, r *http.Request) {
// 	io.WriteString(w, "Hello World!")
// }

func main() {
	router := gin.Default()
	//tmpl読み込み
	router.LoadHTMLGlob("view/index.html")

	port := os.Getenv("PORT")
	fmt.Println("port=", port)
	// http.HandleFunc("/", hello)
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "index.html", gin.H{})
	})
	http.ListenAndServe(":"+port, nil)
}
