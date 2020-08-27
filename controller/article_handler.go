package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	DB *sqlx.DB
}

//TOPページ　状態が死んでいるスレッドの取得は行わない
func HelloPage(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{})
}

func PostPage(c *gin.Context) {
	c.Request.ParseForm()
	name := c.Request.Form["name"][0]
	c.HTML(200, "index0.html", gin.H{
		"name": name,
	})
}
