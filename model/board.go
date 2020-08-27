package model

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type BoardDB struct {
	ID       string    `db:"id"`
	UserName string    `db:"name"`
	Text     string    `db:"text"`
	Time     time.Time `db:"time"`
}

//投稿時にTABLEを作成する
func CreateBoard(db *sqlx.DB, boardName string) error {
	fmt.Println("CreateBoard")
	query := `CREATE TABLE ` + boardName + ` (
		id SERIAL NOT NULL,
		name VARCHAR(10),
		text VARCHAR(200),
		time TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP
		)`
	stmt, err := db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		fmt.Println("STMT NG")
		log.Fatal(err)
	}
	fmt.Println("STMT OK")
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println("STMT NG2")
		log.Fatal(err)
	}
	fmt.Println("STMT OK2")
	return nil
}

//作成したTABLEにデータの追加
//idはTABLE検索用
func InsertBoard(db *sqlx.DB, boardName, userName, text string) error {
	fmt.Println("InsertBoard")
	query := "INSERT INTO " + boardName + " (name,text) VALUES ($1,$2)"
	stmt, err := db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(userName, text)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

//スレッドの内容を見る処理
func SeeBoardPages(db *sqlx.DB, boardName string) ([]BoardDB, error) {
	fmt.Println("SeeBoardPages")
	result := make([]BoardDB, 0)
	//SELECTを実行。db.Queryの代わりにdb.Queryxを使う。
	query := "SELECT * FROM " + boardName
	rows, err := db.Queryx(query)
	if err != nil {
		log.Fatal(err)
	}
	var DB BoardDB
	for rows.Next() {
		//rows.Scanの代わりにrows.StructScanを使う
		err := rows.StructScan(&DB)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, DB)
	}
	return result, err
}

func DeleteBoardWrites(db *sqlx.DB, boardName, id string) error {
	fmt.Println("DeleteBoardWrites")
	query := "UPDATE " + boardName + " SET text = $1 WHERE id = $2"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("この書き込みは削除されました", id)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
