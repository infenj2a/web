package util

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

func InitDB() *sqlx.DB {
	// 接続情報は heroku config :get pacific-bayou-32131 で取得
	db, err := sqlx.Open("postgres", "host=ec2-54-91-178-234.compute-1.amazonaws.com user=zsendmswafsfvc dbname=d304v44t1lic6h sslmode=require password=f8286a9a23f6af42df5aa2ff7d80e3432f8e073c3ccb454c1ada162d5e45831a")
	// db, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	//defer db.Close()を付与するとInitDBが終わった際に閉じてしまうので付与していない
	if err != nil {
		log.Fatal(err)
	}
	return db
}
