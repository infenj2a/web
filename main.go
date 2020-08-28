// main.go

package main

import (
	"github.com/gin-gonic/gin"
	melody "gopkg.in/olahol/melody.v1"
	"main/controller"
	"main/util"
	"net/http"
	"os"
)

func main() {
	server := controller.Server{
		DB: util.InitDB(),
	}

	r := gin.Default()
	m := melody.New()

	//css読み込み準備
	r.Static("/chat_css", "./view/css_chat")
	r.Static("/page_css", "./view/css_page")
	//tmpl読み込み
	r.LoadHTMLGlob("view/*tmpl")

	//GET
	r.GET("/", server.GetArticlePage)
	r.GET("/page/:page", server.GetArticlePage)
	r.GET("/all", server.AllGetArticlePage)
	r.GET("/new", server.GetArticleNewPage)
	r.GET("/board/:id", server.SeeBoardPage)

	//POST
	r.POST("/search", server.PostSearchPage)
	r.POST("/new", server.PostArticleNewPage)
	r.POST("/write/:boardName", server.WriteBoardPage)
	r.POST("/cancel/:boardName/:id", server.DeleteBoardWrite)
	r.POST("/delete/:id", server.DeleteArticlePage)
	r.POST("/drop/:id/:status", server.DropArticlePage)
	r.POST("/status/:id/:status", server.PostStatusChange)

	//WebSocket
	r.GET("/lobby", server.Lobby)
	r.GET("/room/:name", func(c *gin.Context) {
		c.HTML(http.StatusOK, "room.tmpl", gin.H{
			"Name": c.Param("name"),
		})
	})

	r.GET("/room/:name/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.BroadcastFilter(msg, func(q *melody.Session) bool {
			return q.Request.URL.Path == s.Request.URL.Path
		})
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
