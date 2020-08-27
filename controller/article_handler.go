package controller

import (
	"github.com/gin-gonic/gin"
)

//TOPページ　状態が死んでいるスレッドの取得は行わない
func HelloPage(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{})
}

func PostPage(c *gin.Context) {
	name := c.PostForm("name")
	c.HTML(200, "index0.html", gin.H{
		"name": name,
	})
}
