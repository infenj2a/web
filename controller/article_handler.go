package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"main/model"
	"net/http"
)

type Server struct {
	DB *sqlx.DB
}

func (s *Server) GetPage(c *gin.Context) {
	errMsg := ""
	articles, err := model.GetArticle(s.DB)
	if err != nil {
		errMsg = "エラー発生"
		articles = []model.ArticleDB{}
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	c.HTML(200, "index.html", gin.H{
		"articles": articles,
		"errMsg":   &errMsg,
	})
	return
}

func (s *Server) PostPage(c *gin.Context) {
	errMsg := ""
	articles, err := model.PostArticle(s.DB)
	if err != nil {
		errMsg = "エラー発生"
		articles = []model.ArticleDB{}
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	c.HTML(200, "index.html", gin.H{
		"articles": articles,
		"errMsg":   &errMsg,
	})
	return
}
