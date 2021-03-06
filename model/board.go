package model

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type BoardDB struct {
	ID     string    `db:"id"`
	Userid string    `db:"userid"`
	Name   string    `db:"name"`
	Text   string    `db:"text"`
	Time   time.Time `db:"time"`
	Time2  string
}

//投稿時にTABLEを作成する
func CreateBoard(db *sqlx.DB, boardName string) error {
	fmt.Println("CreateBoard")
	query := `CREATE TABLE ` + boardName + ` (
		id SERIAL NOT NULL,
		userid VARCHAR(10),
		name VARCHAR(10),
		text VARCHAR(200),
		time TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP
		)`
	stmt, err := db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

//作成したTABLEにデータの追加
//idはTABLE検索用
func InsertBoard(db *sqlx.DB, boardName, userid, name, text string) error {
	fmt.Println("InsertBoard")
	fmt.Println("userid=", userid, "name=", name)
	query := "INSERT INTO " + boardName + " (userid,name,text) VALUES ($1,$2,$3)"
	stmt, err := db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(userid, name, text)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

//スレッドの内容を見る処理
func SeeBoardPages(db *sqlx.DB, boardName string) ([]BoardDB, error) {
	fmt.Println("SeeBoardPages")
	result := make([]BoardDB, 0)
	query := "SELECT * FROM " + boardName + " order by id asc"
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
		// 時間の整形
		DB.Time2 = StringTime(DB.Time)
		result = append(result, DB)
	}
	return result, err
}

func DeleteBoardWrites(db *sqlx.DB, boardName, id string) error {
	fmt.Println("DeleteBoardWrites")
	query := "UPDATE " + boardName + " SET text = $1 WHERE id = $2"
	stmt, err := db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("この書き込みは削除されました", id)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
