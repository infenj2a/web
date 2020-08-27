package model

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"sort"
	"time"
)

type ArticleDB struct {
	ID         string    `db:"id"`
	Title      string    `db:"title"`
	Status     string    `db:"status"`
	Updatetime time.Time `db:"updatetime"`
}

type CountPage struct {
	Page int
}

func GetArticles(db *sqlx.DB, index int) ([]ArticleDB, int, []CountPage, error) {
	fmt.Println("GetArticles")
	page := 0
	if index != 1 {
		page += (index - 1) * 5
	}
	resultStruct := make([]ArticleDB, 0)
	//SELECTを実行。db.Queryの代わりにdb.Queryxを使う。
	rows, err := db.Queryx(`SELECT * FROM articles WHERE status = ? ORDER BY updatetime DESC LIMIT 5 OFFSET ?;`, "Alive", page)
	if err != nil {
		log.Fatal(err)
	}
	var DB ArticleDB
	for rows.Next() {
		//rows.Scanの代わりにrows.StructScanを使う
		err := rows.StructScan(&DB)
		if err != nil {
			fmt.Println("err")
			log.Fatal(err)
			return nil, 0, nil, err
		}
		resultStruct = append(resultStruct, DB)
	}
	rows, err = db.Queryx(`SELECT COUNT(status) FROM articles WHERE status = "Alive"`)
	resultCount := checkCount(rows)
	indexStruct, div := exprCount(resultCount, index)
	if err != nil {
		log.Fatal(err)
	}
	return resultStruct, div, indexStruct, err
}

func checkCount(rows *sqlx.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
	}
	return count
}

//index は基準
func exprCount(count, index int) ([]CountPage, int) {
	fmt.Println("exprCount")
	fmt.Println("count=", count)
	fmt.Println("index=", index)
	indexStruct := make([]CountPage, 0)
	div := (count % 5)
	if div == 0 {
		div = (count / 5)
	} else {
		div = (count / 5) + 1
	}
	if count <= 5 || div < index {
		fmt.Println("5件以内or異常ページ感知")
		return indexStruct, div
	} else {
		fmt.Println("for分開始")
		slice := []int{index}
		for x := 1; x < 5; x++ {
			if z := index + x; div >= z {
				slice = append(slice, z)
			}
			if z := index - x; 0 < z {
				slice = append(slice, z)
			}
			if len(slice) > 4 {
				break
			}
		}
		fmt.Println("ソート前スライス=", slice)
		sort.Sort(sort.IntSlice(slice))
		fmt.Println("ソート後スライス=", slice)
		for _, v := range slice {
			indexStruct = append(indexStruct, CountPage{v})
		}
		return indexStruct, div
	}
}
func AllGetArticles(db *sqlx.DB) ([]ArticleDB, error) {
	fmt.Println("AllGetArticles")
	result := make([]ArticleDB, 0)
	//SELECTを実行。db.Queryの代わりにdb.Queryxを使う。
	rows, err := db.Queryx(`SELECT * FROM articles;`)
	if err != nil {
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
		result = append(result, DB)
	}
	return result, err
}

func GetArticleOne(db *sqlx.DB, id string) ([]ArticleDB, error) {
	fmt.Println("GetArticles")
	result := make([]ArticleDB, 0)
	//SELECTを実行。db.Queryの代わりにdb.Queryxを使う。
	rows, err := db.Queryx(`SELECT * FROM articles WHERE ID = ?;`, id)
	if err != nil {
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
		result = append(result, DB)
	}
	return result, err
}

// NewArticlePosts = INSERT
func PostArticleNewPages(db *sqlx.DB, title string) int {
	fmt.Println("PostArticleNewPages")
	stmt, err := db.Prepare(`INSERT INTO articles (title, status) VALUES (?, ?);`)
	if err != nil {
		fmt.Println("stmt err")
		log.Fatal(err)
		return 0
	}
	resp, err := stmt.Exec(title, "Alive")
	if err != nil {
		fmt.Println("resp err")
		log.Fatal(err)
		return 0
	}
	//影響のあった件数が返る index + errが返る
	// id, _ := resp.RowsAffected()
	//一番最後のIDを取得したい場合 index + errが返る
	id, err := resp.LastInsertId()
	if err != nil {
		fmt.Println("LastInsertId err")
		return 0
	}
	return int(id)
}

func DeleteArticlePages(db *sqlx.DB, id int) error {
	fmt.Println("DeleteArticlePages")
	//SELECTを実行。db.Queryの代わりにdb.Queryxを使う。
	stmt, err := db.Prepare(`UPDATE articles SET status = ? WHERE ID = ?;`)
	if err != nil {
		fmt.Println("Delete UPDATE err")
		log.Fatal(err)
		return err
	}
	_, err = stmt.Exec("Dead", id)
	if err != nil {
		fmt.Println("Delete EXEC err")
		log.Fatal(err)
		return err
	}
	return nil
}

func DropArticlePages(db *sqlx.DB, id int, status string) error {
	fmt.Println("DropArticlePages")
	if status == "Dead" {
		stmt, err := db.Prepare(`DELETE FROM articles WHERE ID = ?;`)
		if err != nil {
			fmt.Println("Drop stmt err")
			return err
		}
		_, err = stmt.Exec(id)
		if err != nil {
			fmt.Println("Drop EXEC err")
			return err
		}
	}
	return nil
}

func StatusChangePage(db *sqlx.DB, id int, status string) error {
	fmt.Println("StatusChangePage")
	switch status {
	case "Alive":
		status = "Dead"
	case "Dead":
		status = "Alive"
	default:
		status = "不明"
	}
	//SELECTを実行。db.Queryの代わりにdb.Queryxを使う。
	stmt, err := db.Prepare(`UPDATE articles SET status = ? WHERE ID = ?;`)
	if err != nil {
		fmt.Println("StatusChange stmt err")
		log.Fatal(err)
		return err
	}
	_, err = stmt.Exec(status, id)
	if err != nil {
		fmt.Println("StatusChange EXEC err")
		log.Fatal(err)
		return err
	}
	return nil
}

func UpdateArticleTimes(db *sqlx.DB, id string) error {
	fmt.Println("UpdateArticleTimes")
	//SELECTを実行。db.Queryの代わりにdb.Queryxを使う。
	stmt, err := db.Prepare(`UPDATE articles SET updatetime = ? WHERE ID = ?;`)
	if err != nil {
		fmt.Println("UPDATETIME stmt err")
		log.Fatal(err)
		return err
	}
	_, err = stmt.Exec(time.Now(), id)
	if err != nil {
		fmt.Println("UPDATETIME EXEC err")
		log.Fatal(err)
		return err
	}
	return nil
}

func PostSearchPages(db *sqlx.DB, title string) ([]ArticleDB, error) {
	fmt.Println("PostSearchPages")
	title = "%" + title + "%"
	result := make([]ArticleDB, 0)
	//SELECTを実行。db.Queryの代わりにdb.Queryxを使う。
	rows, err := db.Queryx(`SELECT * FROM articles WHERE title LIKE ?;`, title)
	if err != nil {
		fmt.Println("PostSearchPages SELECT err")
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
		result = append(result, DB)
	}
	return result, err
}
