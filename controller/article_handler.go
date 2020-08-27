package controller

import (
	"github.com/gin-gonic/gin"
)

//TOPページ　状態が死んでいるスレッドの取得は行わない
func HelloPage(ctx *gin.Context) {
	ctx.HTML(200, "index.html", gin.H{})
}

func PostPage(ctx *gin.Context) {
	ctx.Request.ParseForm()
	name := ctx.Request.Form["name"][0]
	ctx.HTML(200, "index0.html", gin.H{
		"name": name,
	})
}
