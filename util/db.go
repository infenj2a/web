package util

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

func InitDB() *sqlx.DB {
	// 接続情報は heroku config :get リポジトリの名前 で取得可能
	// デフォルトsslmodeはdisable
	db, err := sqlx.Open("postgres", "host= user= dbname= sslmode=require password=")
	// 上記の記載が面倒な場合は下記で実行が可能
	// db, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))

	// defer db.Close()を付与するとInitDBが終わった際に閉じてしまうので付与していない
	// 接続数に限りがあるので将来性を考えるなら考慮するべきところ
	if err != nil {
		log.Fatal(err)
	}
	// articlesテーブルが無ければ作成
	query := `CREATE TABLE IF NOT EXISTS articles(
		id SERIAL NOT NULL,
		title VARCHAR(40),
		status VARCHAR(10),
		updatetime TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
		)`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
