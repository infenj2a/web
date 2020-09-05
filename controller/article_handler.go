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

type UserDB struct {
	UserID string `db:"userid"`
	userPW string `db:"password"`
}

var userSession UserDB
var errMsg string

func DeleteSession() {
	userSession.UserID = ""
}

// ログインページ
func (s *Server) LoginPage(c *gin.Context) {
	fmt.Println("LoginPage")
	// セッション情報が保持されている場合は/homeへ移動
	if userSession.UserID == "" {
		fmt.Println("Session NG")
		c.HTML(http.StatusOK, "login.tmpl", gin.H{})
		return
	}
	fmt.Println("Session OK")
	c.Redirect(http.StatusMovedPermanently, "/top")
}

// ログインページよりID,PWが入力された際の処理
func (s *Server) LoginUser(c *gin.Context) {
	fmt.Println("LoginUser")
	// フォーム情報の取得
	c.Request.ParseForm()
	userSession.UserID = c.Request.Form["ID"][0]
	userSession.userPW = c.Request.Form["PW"][0]
	// フォーム情報確認
	if userSession.UserID == "" {
		errMsg = "IDを入力してください"
		// セッション情報の破棄
		DeleteSession()
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"errMsg": &errMsg,
		})
		errMsg = ""
		return
	}
	if userSession.userPW == "" {
		errMsg = "PWを入力してください"
		// セッション情報の破棄
		DeleteSession()
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"errMsg": &errMsg,
		})
		errMsg = ""
		return
	}
	// ログインが完了した場合、/homeページへ移動
	err := model.LoginUsers(s.DB, userSession.UserID, userSession.userPW)
	if err != 1 {
		//ID or PWが間違っている
		errMsg = "IDまたはパスワードが間違っています。"
		// セッション情報の破棄
		DeleteSession()
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"errMsg": &errMsg,
		})
		errMsg = ""
		return
	}
	// 正常ログイン
	c.Redirect(http.StatusMovedPermanently, "/top")
}

// ログインページより新規登録ページ
func (s *Server) CreatePage(c *gin.Context) {
	fmt.Println("userSession=", userSession.UserID)
	// セッション情報の破棄
	DeleteSession()
	fmt.Println("userSession=", userSession.UserID)
	c.HTML(http.StatusOK, "create.tmpl", gin.H{})
	errMsg = ""
	return
}

// 新規登録ページより入力された時の処理
func (s *Server) CreateUser(c *gin.Context) {
	fmt.Println("CreateUser")
	// フォーム情報の取得
	c.Request.ParseForm()
	userSession.UserID = c.Request.Form["ID"][0]
	userSession.userPW = c.Request.Form["PW"][0]
	// フォーム情報確認
	if userSession.UserID == "" {
		errMsg = "IDを入力してください"
		// セッション情報の破棄
		DeleteSession()
		c.HTML(http.StatusOK, "create.tmpl", gin.H{
			"errMsg": &errMsg,
		})
		errMsg = ""
		return
	}
	if userSession.userPW == "" {
		errMsg = "パスワードを入力してください"
		// セッション情報の破棄
		DeleteSession()
		c.HTML(http.StatusOK, "create.tmpl", gin.H{
			"errMsg": &errMsg,
		})
		errMsg = ""
		return
	}
	// 新規登録が完了した場合、/homeページへ移動
	err := model.CreateUsers(s.DB, userSession.UserID, userSession.userPW)
	if err != nil {
		// すでに名前が他のユーザーに登録されていた場合
		errMsg = "申し訳ありません。既に使用されている名前ですので、他の名前をご入力ください。"
		// セッション情報の破棄
		DeleteSession()
		c.HTML(http.StatusOK, "create.tmpl", gin.H{
			"errMsg": &errMsg,
		})
		errMsg = ""
		return
	}
	// 正常ログイン
	errMsg = ""
	c.Redirect(http.StatusMovedPermanently, "/top")
}

func (s *Server) Logout(c *gin.Context) {
	fmt.Println("LogoutUser")
	// セッション情報の破棄
	DeleteSession()
	errMsg = ""
	c.HTML(http.StatusOK, "login.tmpl", gin.H{})
}

// /homeページ statusが Deadのレコードは取得せず、5件ずつ取得表示
func (s *Server) GetArticlePage(c *gin.Context) {
	index := c.Param("page")
	// ページャー用page
	page, err := strconv.Atoi(index)
	if err != nil {
		page = 1
	}
	//articles 		= ページ情報
	//div 				= ページャー用 計算用
	//pageStruct 	= ページャー用 配列
	articles, div, pageStruct, err := model.GetArticles(s.DB, page)
	if err != nil {
		articles = []model.ArticleDB{}
		pageStruct = []model.CountPage{}
		errMsg = "エラーが発生しました。ページを更新して下さい。"
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
		"Session":    userSession.UserID,
	})
	errMsg = ""
	return
}

//スレッド一覧の全取得
func (s *Server) AllGetArticlePage(c *gin.Context) {
	errMsg := ""
	articles, err := model.AllGetArticles(s.DB)
	if err != nil {
		errMsg = "エラー発生"
		articles = []model.ArticleDB{}
		c.Redirect(http.StatusMovedPermanently, "/top")
		return
	}
	c.HTML(http.StatusOK, "all.tmpl", gin.H{
		"articles": articles,
		"errMsg":   &errMsg,
		"Session":  userSession.UserID,
	})
	errMsg = ""
	return
}

//新規スレッド作成
func (s *Server) GetArticleNewPage(c *gin.Context) {
	c.HTML(http.StatusOK, "new.tmpl", gin.H{
		"Session": userSession.UserID,
	})
	errMsg = ""
	return
}

//新規スレッド作成要求が発生した場合
func (s *Server) PostArticleNewPage(c *gin.Context) {
	// フォーム情報の取得
	c.Request.ParseForm()
	article := new(model.ArticleDB)
	board := new(model.BoardDB)
	article.Title = c.Request.Form["title"][0]
	board.Name = c.Request.Form["name"][0]
	board.Text = c.Request.Form["text"][0]
	// フォーム情報確認
	if article.Title == "" {
		c.HTML(http.StatusBadRequest, "new.tmpl", gin.H{
			"errMsg":  "タイトルを入力して下さい",
			"Session": userSession.UserID,
		})
		errMsg = ""
		return
	}
	if board.Text == "" {
		c.HTML(http.StatusBadRequest, "new.tmpl", gin.H{
			"errMsg":  "テキストが空です",
			"Session": userSession.UserID,
		})
		errMsg = ""
		return
	}
	if board.Name == "" {
		board.Name = "名無しさん"
	}
	//一度、ArticleDBに登録を実行し、idを取得
	id := model.PostArticleNewPages(s.DB, userSession.UserID, article.Title)
	if id == 0 {
		c.HTML(http.StatusBadRequest, "new.tmpl", gin.H{
			"errMsg":  "登録できませんでした",
			"Session": userSession.UserID,
		})
		errMsg = ""
		return
	}
	//スレッド専用のTABLEを作成 id名のスレッドを作成
	tableID := strconv.Itoa(id)
	boardName := "board" + tableID
	if err := model.CreateBoard(s.DB, boardName); err != nil {
		c.HTML(http.StatusBadRequest, "new.tmpl", gin.H{
			"errMsg":  "専用スレッドの作成に失敗しました",
			"Session": userSession.UserID,
		})
		errMsg = ""
		return
	}
	//作成したBoard専用のDBにArticleDBと同じ内容を追記(タイトルは格納しない)
	if err := model.InsertBoard(s.DB, boardName, userSession.UserID, board.Name, board.Text); err != nil {
		c.HTML(http.StatusBadRequest, "new.tmpl", gin.H{
			"errMsg":  "専用スレッドへの書き込みに失敗しました",
			"Session": userSession.UserID,
		})
		errMsg = ""
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
		c.Redirect(http.StatusMovedPermanently, "/top")
		return
	}

	//専用スレッドの情報を取得
	board, err := model.SeeBoardPages(s.DB, boardName)
	if err != nil {
		errMsg = "スレッドのロード失敗"
		c.Redirect(http.StatusMovedPermanently, "/top")
		return
	}
	boardTitle := article[0].Title
	c.HTML(http.StatusOK, "board.tmpl", gin.H{
		"boardTitle": boardTitle,
		"board":      board,
		"boardName":  boardName,
		"errMsg":     &errMsg,
		"Session":    userSession.UserID,
	})
	errMsg = ""
	return
}

//専用スレッド内からの書き込み追加
func (s *Server) WriteBoardPage(c *gin.Context) {
	fmt.Println("WriteBoardPage")
	c.Request.ParseForm()
	boardName := c.Param("boardName")
	board := new(model.BoardDB)
	board.Name = c.Request.Form["name"][0]
	board.Text = c.Request.Form["text"][0]
	id := strings.Replace(boardName, "board", "", 1)
	if board.Text == "" {
		c.Redirect(http.StatusMovedPermanently, "/board/"+id)
		return
	}
	if board.Name == "" {
		board.Name = "名無しさん"
	}

	//ArticleDBの対象スレッドの時間を更新
	err := model.UpdateArticleTimes(s.DB, id)
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/board/"+id)
		return
	}

	//BoradDBにINSERT
	err = model.InsertBoard(s.DB, boardName, userSession.UserID, board.Name, board.Text)
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
		c.Redirect(http.StatusMovedPermanently, "/top")
		return
	}
	err = model.DeleteArticlePages(s.DB, id)
	if err != nil {
		log.Fatal(err)
	}
	c.Redirect(http.StatusMovedPermanently, "/top")
	return
}

//全体表示から完全削除の要求があった場合
func (s *Server) DropArticlePage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/top")
		return
	}
	status := c.Param("status")
	err = model.DropArticlePages(s.DB, id, status)
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/top")
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/all")
	return
}

//全体表示からステータスの変更要求があった場合
func (s *Server) PostStatusChange(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/top")
		return
	}
	status := c.Param("status")
	err = model.StatusChangePage(s.DB, id, status)
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/top")
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
			"Session":  userSession.UserID,
		})
		errMsg = ""
		return
	}
	articles, err := model.PostSearchPages(s.DB, article.Title)
	if err != nil {
		errMsg = "エラーが発生しました"
		articles = []model.ArticleDB{}
		c.Redirect(http.StatusMovedPermanently, "/top")
		return
	}
	if len(articles) == 0 {
		errMsg = "検索結果が1件もありません"
	}
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"articles": articles,
		"errMsg":   &errMsg,
		"Session":  userSession.UserID,
	})
	errMsg = ""
	return
}

func (s *Server) Lobby(c *gin.Context) {
	c.HTML(http.StatusOK, "lobby.tmpl", gin.H{
		"Session": userSession.UserID,
	})
	errMsg = ""
	return
}

func (s *Server) Room(c *gin.Context) {
	c.HTML(http.StatusOK, "room.tmpl", gin.H{
		"Name":    c.Param("name"),
		"Session": userSession.UserID,
	})
	errMsg = ""
	return
}
