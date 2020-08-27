package model

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type ArticleDB struct {
	ID         string    `db:"id"`
	Title      string    `db:"title"`
	Status     string    `db:"status"`
	Updatetime time.Time `db:"updatetime"`
}

func GetArticle(db *sqlx.DB) ([]ArticleDB, error) {
	fmt.Println("GetArticle")
	resultStruct := make([]ArticleDB, 0)
	rows, err := db.Queryx("SELECT * FROM articles")
	if err != nil {
		fmt.Println(err)
		fmt.Println(rows)
		fmt.Println("SELECTエラー")
		log.Fatal(err)
	}
	var DB ArticleDB
	for rows.Next() {
		//rows.Scanの代わりにrows.StructScanを使う
		err := rows.StructScan(&DB)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		resultStruct = append(resultStruct, DB)
	}
	return resultStruct, nil
}

func PostArticle(db *sqlx.DB) ([]ArticleDB, error) {
	fmt.Println("PostArticle")
	resultStruct := make([]ArticleDB, 0)
	rows, err := db.Queryx(`SELECT * FROM articles WHERE id = 2;`)
	if err != nil {
		fmt.Println("SELECTエラー")
		log.Fatal(err)
	}
	var DB ArticleDB
	for rows.Next() {
		//rows.Scanの代わりにrows.StructScanを使う
		err := rows.StructScan(&DB)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		resultStruct = append(resultStruct, DB)
	}
	return resultStruct, nil
}
