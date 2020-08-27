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
	_, err := db.Query(`CREATE TABLE ` + boardName + ` (
		id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
		name VARCHAR(10),
		text VARCHAR(200),
		time DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
		);`)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

//作成したTABLEにデータの追加
//idはTABLE検索用
func InsertBoard(db *sqlx.DB, boardName, userName, text string) error {
	fmt.Println("InsertBoard")
	stmt, err := db.Prepare(`INSERT INTO ` + boardName + ` (name,text) VALUES (?,?);`)
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = stmt.Exec(userName, text)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

//スレッドの内容を見る処理
func SeeBoardPages(db *sqlx.DB, boardName string) ([]BoardDB, error) {
	fmt.Println("SeeBoardPages")
	result := make([]BoardDB, 0)
	//SELECTを実行。db.Queryの代わりにdb.Queryxを使う。
	rows, err := db.Queryx(`SELECT * FROM ` + boardName + `;`)
	if err != nil {
		log.Fatal(err)
	}
	var DB BoardDB
	for rows.Next() {
		//rows.Scanの代わりにrows.StructScanを使う
		err := rows.StructScan(&DB)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		result = append(result, DB)
		fmt.Println(result)
	}
	return result, err
}

func DeleteBoardWrites(db *sqlx.DB, boardName, id string) error {
	fmt.Println("DeleteBoardWrites")
	stmt, err := db.Prepare(`UPDATE ` + boardName + ` SET text = (?) WHERE ID = ?;`)
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = stmt.Exec("この書き込みは削除されました", id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("正常")
	return nil
}
