package model

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

//ID,PWが入力された際の処理
func LoginUsers(db *sqlx.DB, id, pw string) int {
	fmt.Println("LoginUsers")
	query := "SELECT * FROM userdate WHERE userid = $1 AND password = $2"
	rows, err := db.Queryx(query, id, pw)
	if err != nil {
		fmt.Println("SQL文エラー")
		return 0
	}
	// 1件だけ取得できたらログイン
	i := 0
	for rows.Next() {
		i += 1
	}
	if i == 1 {
		return 1
	}
	return 0
}

//ID,PWが入力された際の処理
func CreateUsers(db *sqlx.DB, id, pw string) error {
	fmt.Println("CreateUsers")
	query := "INSERT INTO userdate (userid,password) VALUES($1,$2)"
	stmt, err := db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		fmt.Println("SQL文エラー")
		return err
	}
	_, err = stmt.Exec(id, pw)
	if err != nil {
		fmt.Println("useridが同じものは登録はできない")
		return err
	}
	return nil
}
