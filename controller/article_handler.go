package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"main/model"
	"net/http"
)

type Server struct {
	DB *sqlx.DB
}

func (s *Server) GetPage(c *gin.Context) {
	fmt.Println("GetPage")
	errMsg := ""
	articles, err := model.GetArticle(s.DB)
	if err != nil {
		fmt.Println("エラー発生")
		errMsg = "エラー発生"
		articles = []model.ArticleDB{}
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	fmt.Println("取得OK")
	c.HTML(200, "index.html", gin.H{
		"articles": articles,
		"errMsg":   &errMsg,
	})
	return
}

func (s *Server) PostPage(c *gin.Context) {
	fmt.Println("PostPage")
	errMsg := ""
	articles, err := model.PostArticle(s.DB)
	if err != nil {
		fmt.Println("エラー発生")
		errMsg = "エラー発生"
		articles = []model.ArticleDB{}
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	fmt.Println("取得OK")
	c.HTML(200, "index.html", gin.H{
		"articles": articles,
		"errMsg":   &errMsg,
	})
	return
}
