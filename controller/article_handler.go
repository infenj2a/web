package controller

import (
	"github.com/gin-gonic/gin"
)

//TOPページ　状態が死んでいるスレッドの取得は行わない
func HelloPage(ctx *gin.Context) {
	ctx.HTML(200, "index.html", gin.H{})
}
