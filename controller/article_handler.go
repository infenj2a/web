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

//TOPページ statusがDeadのレコードは取得せず、getRecord件数単位で取得
func (s *Server) GetArticlePage(c *gin.Context) {
	index := c.Param("page")
	page, err := strconv.Atoi(index)
	if err != nil {
		page = 1
	}
	errMsg := ""
	//articles 		= ページ情報
	//div 				= ページャー用 計算用
	//pageStruct 	= ページャー用 配列
	articles, div, pageStruct, err := model.GetArticles(s.DB, page)
	if err != nil {
		articles = []model.ArticleDB{}
		pageStruct = []model.CountPage{}
		errMsg = "エラー発生"
	}
	//prevPage = ページャー用 戻る判定 0の時はindex.htmlで表示されない仕様
	//nextPage = ページャー用 進む判定 0の時はindex.htmlで表示されない仕様
	prevPage := page - 1
	nextPage := page + 1
	if nextPage > div {
		nextPage = 0
	}
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

//新規スレッド作成要求が発生した場合
func (s *Server) PostArticleNewPage(c *gin.Context) {
	//POSTされた情報を分解
	c.Request.ParseForm()
	article := new(model.ArticleDB)
	board := new(model.BoardDB)
	article.Title = c.Request.Form["title"][0]
	board.UserName = c.Request.Form["name"][0]
	board.Text = c.Request.Form["text"][0]

	//情報に不備がないかを確認
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

	//一度、ArticleDBに登録を実行し、挿入したidを取得
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

	//作成したBoard専用のDBにArticleDBと同じ内容を追記(タイトルは格納しない)
	if err := model.InsertBoard(s.DB, boardName, board.UserName, board.Text); err != nil {
		c.HTML(http.StatusBadRequest, "new.tmpl", gin.H{
			"errMsg": "専用スレッドへの書き込みに失敗しました",
		})
		return
	}
	//新規作成スレッドの表示
	c.Redirect(http.StatusMovedPermanently, "/board/"+tableID)
	return
}

//TOPページより特定のスレッドを開いた際の処理
func (s *Server) SeeBoardPage(c *gin.Context) {
	fmt.Println("SeeBoardPage")
	id := c.Param("id")
	boardName := "board" + id
	errMsg := ""

	//idよりスレッドのタイトルの取得
	article, err := model.GetArticleOne(s.DB, id)
	if err != nil {
		errMsg = "スレッドのタイトルが正常に取得できませんでした"
		article = []model.ArticleDB{}
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}

	//専用スレッドの情報を取得
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

//専用スレッド内からの書き込み追加
func (s *Server) WriteBoardPage(c *gin.Context) {
	fmt.Println("WriteBoardPage")
	c.Request.ParseForm()
	boardName := c.Param("boardName")
	board := new(model.BoardDB)
	board.UserName = c.Request.Form["name"][0]
	board.Text = c.Request.Form["text"][0]
	id := strings.Replace(boardName, "board", "", 1)
	if board.Text == "" {
		c.Redirect(http.StatusMovedPermanently, "/board/"+id)

	}
	if board.UserName == "" {
		board.UserName = "名無しさん"
	}

	//ArticleDBの対象スレッドの時間を更新
	err := model.UpdateArticleTimes(s.DB, id)
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/board/"+id)
		return
	}

	//BoradDBにINSERT
	err = model.InsertBoard(s.DB, boardName, board.UserName, board.Text)
	if err != nil {
		log.Fatal(err)
	}
	c.Redirect(http.StatusMovedPermanently, "/board/"+id)
	return
}

//専用スレッドより書き込み削除が要求された場合
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

//削除
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

//全体表示から完全削除の要求があった場合
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

//全体表示からステータスの変更要求があった場合
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

//TOPページからの検索
func (s *Server) PostSearchPage(c *gin.Context) {
	c.Request.ParseForm()
	errMsg := ""
	article := new(model.ArticleDB)
	article.Title = c.Request.Form["search"][0]
	if article.Title == "" {
		errMsg = "検索キーワードを入力してください"
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"articles": []model.ArticleDB{},
			"errMsg":   &errMsg,
		})
		return
	}
	articles, err := model.PostSearchPages(s.DB, article.Title)
	if err != nil {
		errMsg = "エラーが発生しました"
		articles = []model.ArticleDB{}
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	if len(articles) == 0 {
		errMsg = "検索結果が1件もありません"
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
