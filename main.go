package main

import (
	"github.com/gin-gonic/gin"
	melody "gopkg.in/olahol/melody.v1"
	"main/controller"
	"main/util"
	"os"
)

func main() {
	server := controller.Server{
		DB: util.InitDB(),
	}

	r := gin.Default()
	m := melody.New()

	//css読み込み準備
	r.Static("/css_chat", "./view/css_chat")
	r.Static("/css_page", "./view/css_page")
	r.Static("/css_login", "./view/css_login")

	//tmpl読み込み
	r.LoadHTMLGlob("view/*tmpl")

	// ログインページ
	r.GET("/", server.LoginPage)
	// ログインページ(topからログアウトが正常に動作しないのでPOSTで処理)
	r.POST("/", server.LoginPage)
	// ログインページより、ログイン実行
	r.POST("/login", server.LoginUser)
	// 新規登録ページ
	r.GET("/create", server.CreatePage)
	// 新規登録ページより、新規登録実行
	r.POST("/create", server.CreateUser)
	// ゲストログインは直接/homeへ移動
	r.POST("/bye", server.Logout)

	// ホームページ5件ずつ出力を実行
	r.GET("/top", server.GetArticlePage)
	// ホームページよりページャー機能で他のページに飛んだ時の処理
	r.GET("/page/:page", server.GetArticlePage)
	//　ホームから検索機能を使用した際の処理
	r.POST("/search", server.PostSearchPage)

	// 全件表示ページ
	r.GET("/all", server.AllGetArticlePage)
	// 新規投稿ページ
	r.GET("/new", server.GetArticleNewPage)
	// 新規投稿ページより投稿が実行された際の処理
	r.POST("/new", server.PostArticleNewPage)

	// ホームから特定のスレッドをクリック
	r.GET("/board/:id", server.SeeBoardPage)
	// 専用スレッドより新規書き込みがあった際の処理
	r.POST("/write/:boardName", server.WriteBoardPage)
	// 専用スレッドより書き込みの削除依頼が発生した際の処理
	r.POST("/cancel/:boardName/:id", server.DeleteBoardWrite)
	// ホームより削除依頼が発生した際の処理
	r.POST("/delete/:id", server.DeleteArticlePage)
	// 全体表示ページよりステータス変更の依頼が発生した際の処理
	r.POST("/status/:id/:status", server.PostStatusChange)
	// 全体表示ページより削除依頼が発生した際の処理
	r.POST("/drop/:id/:status", server.DropArticlePage)

	// WebSocket
	r.GET("/lobby", server.Lobby)
	r.GET("/room/:name", server.Room)

	// 通信
	r.GET("/room/:name/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	// フィルター
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.BroadcastFilter(msg, func(q *melody.Session) bool {
			return q.Request.URL.Path == s.Request.URL.Path
		})
	})

	//Herokuの環境変数より宛てられたポートを取得
	port := os.Getenv("PORT")
	//ローカル環境用に無い場合は8080で接続するようにする
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
