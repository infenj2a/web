package controller

import (
	"fmt"
	"log"
	"main/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	DB *sqlx.DB
}

//TOPページ　状態が死んでいるスレッドの取得は行わない
func (s *Server) GetArticlePage(c *gin.Context) {
	index := c.Param("page")
	page, err := strconv.Atoi(index)
	if err != nil {
		page = 1
	}
	errMsg := ""
	fmt.Println("page=", page)
	articles, div, pageStruct, err := model.GetArticles(s.DB, page)
	if err != nil {
		articles = []model.ArticleDB{}
		pageStruct = []model.CountPage{}
		errMsg = "エラー発生"
	}
	prevPage := page - 1
	nextPage := page + 1
	fmt.Println("div=", div)
	if nextPage > div {
		nextPage = 0
	}
	fmt.Println("prevPage=", prevPage)
	fmt.Println("nextPage=", nextPage)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"articles":   articles,
		"PrevPage":   prevPage,
		"PageIndex":  page,
		"NextPage":   nextPage,
		"PageStruct": pageStruct,
		"errMsg":     &errMsg,
	})
	return
}

//スレッド一覧の全取得
func (s *Server) AllGetArticlePage(c *gin.Context) {
	errMsg := ""
	articles, err := model.AllGetArticles(s.DB)
	if err != nil {
		errMsg = "エラー発生"
		articles = []model.ArticleDB{}
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	c.HTML(http.StatusOK, "all.tmpl", gin.H{
		"articles": articles,
		"errMsg":   &errMsg,
	})
	return
}

//新規スレッド作成
func (s *Server) GetArticleNewPage(c *gin.Context) {
	c.HTML(http.StatusOK, "new.tmpl", gin.H{})
	return
}

//新規スレッドよりスレッドが作成された場合
func (s *Server) PostArticleNewPage(c *gin.Context) {
	//POSTされた情報を分解
	c.Request.ParseForm()
	article := new(model.ArticleDB)
	board := new(model.BoardDB)
	article.Title = c.Request.Form["title"][0]
	board.UserName = c.Request.Form["name"][0]
	board.Text = c.Request.Form["text"][0]

	if article.Title == "" {
		c.HTML(http.StatusBadRequest, "new.tmpl", gin.H{
			"errMsg": "タイトルを入力して下さい",
		})
		return
	}
	if board.Text == "" {
		c.HTML(http.StatusBadRequest, "new.tmpl", gin.H{
			"errMsg": "テキストが空です",
		})
		return
	}
	if board.UserName == "" {
		board.UserName = "名無しさん"
	}
	//一度、ArticleDBに登録を実行し、idを取得
	id := model.PostArticleNewPages(s.DB, article.Title)
	if id == 0 {
		c.HTML(http.StatusBadRequest, "new.tmpl", gin.H{
			"errMsg": "登録できませんでした",
		})
		return
	}

	//スレッド専用のTABLEを作成 id名のスレッドを作成
	tableID := strconv.Itoa(id)
	boardName := "board" + tableID
	if err := model.CreateBoard(s.DB, boardName); err != nil {
		c.HTML(http.StatusBadRequest, "new.tmpl", gin.H{
			"errMsg": "専用スレッドの作成に失敗しました",
		})
		return
	}
	//作成したBoard専用のDBにArticleDBと同じ内容を追記
	if err := model.InsertBoard(s.DB, boardName, board.UserName, board.Text); err != nil {
		c.HTML(http.StatusBadRequest, "new.tmpl", gin.H{
			"errMsg": "専用スレッドへの書き込みに失敗しました",
		})
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/board/"+tableID)
	return
}

//スレッドの内容を見る
func (s *Server) SeeBoardPage(c *gin.Context) {
	fmt.Println("SeeBoardPage")
	id := c.Param("id")
	boardName := "board" + id
	errMsg := ""
	//スレッドのタイトルを取得しなおす
	article, err := model.GetArticleOne(s.DB, id)
	if err != nil {
		errMsg = "スレッドのタイトルが正常に取得できませんでした"
		article = []model.ArticleDB{}
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}

	//スレッド情報の取得
	board, err := model.SeeBoardPages(s.DB, boardName)
	if err != nil {
		errMsg = "スレッドのロード失敗"
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	boardTitle := article[0].Title
	c.HTML(http.StatusOK, "board.tmpl", gin.H{
		"boardTitle": boardTitle,
		"board":      board,
		"boardName":  boardName,
		"errMsg":     &errMsg,
	})
	return
}

//スレッド内からの新規書き込み
func (s *Server) WriteBoardPage(c *gin.Context) {
	fmt.Println("WriteBoardPage")
	//POSTされた情報を分解
	c.Request.ParseForm()
	boardName := c.Param("boardName")
	board := new(model.BoardDB)
	board.UserName = c.Request.Form["name"][0]
	board.Text = c.Request.Form["text"][0]
	if board.UserName == "" {
		board.UserName = "名無しさん"
	}

	//ArticleDBの対象スレッドの時間を更新
	id := strings.Replace(boardName, "board", "", 1)
	err := model.UpdateArticleTimes(s.DB, id)
	if err != nil {
		log.Fatal(err)
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}

	//BoradDBにINSERT
	if err := model.InsertBoard(s.DB, boardName, board.UserName, board.Text); err != nil {
		log.Fatal(err)
	}
	c.Redirect(http.StatusMovedPermanently, "/board/"+id)
	return
}

func (s *Server) DeleteBoardWrite(c *gin.Context) {
	boardName := c.Param("boardName")
	id := c.Param("id")
	err := model.DeleteBoardWrites(s.DB, boardName, id)
	if err != nil {
		log.Fatal(err)
	}
	id = strings.Replace(boardName, "board", "", 1)
	c.Redirect(http.StatusMovedPermanently, "/board/"+id)
	return
}

func (s *Server) DeleteArticlePage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	fmt.Println(id)
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	err = model.DeleteArticlePages(s.DB, id)
	if err != nil {
		log.Fatal(err)
	}
	c.Redirect(http.StatusMovedPermanently, "/")
	return
}

func (s *Server) DropArticlePage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	status := c.Param("status")
	err = model.DropArticlePages(s.DB, id, status)
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/all")
	return
}

func (s *Server) PostStatusChange(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	status := c.Param("status")
	err = model.StatusChangePage(s.DB, id, status)
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/all")
	return
}

func (s *Server) PostSearchPage(c *gin.Context) {
	c.Request.ParseForm()
	errMsg := ""
	article := new(model.ArticleDB)
	article.Title = c.Request.Form["search"][0]
	if article.Title == "" {
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	articles, err := model.PostSearchPages(s.DB, article.Title)
	if err != nil {
		errMsg = "エラーが発生しました"
		articles = []model.ArticleDB{}
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"articles": articles,
		"errMsg":   &errMsg,
	})
	return
}

func (s *Server) Lobby(c *gin.Context) {
	c.HTML(http.StatusOK, "lobby.tmpl", gin.H{})
}
