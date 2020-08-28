package model

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"sort"
	"time"
)

type ArticleDB struct {
	ID          string    `db:"id"`
	Title       string    `db:"title"`
	Status      string    `db:"status"`
	Updatetime  time.Time `db:"updatetime"`
	Updatetime2 string
}

type CountPage struct {
	Page int
}

const (
	L = "2006-01-02 15:04:05"
)

func GetArticles(db *sqlx.DB, index int) ([]ArticleDB, int, []CountPage, error) {
	fmt.Println("GetArticles")
	page := 0
	if index != 1 {
		page += (index - 1) * 5
	}
	resultStruct := make([]ArticleDB, 0)
	//SELECTを実行。db.Queryの代わりにdb.Queryxを使う。
	query := "SELECT * FROM articles WHERE status = $1 order by updatetime desc LIMIT 5 OFFSET $2"
	rows, err := db.Queryx(query, "Alive", page)
	if err != nil {
		log.Fatal(err)
	}
	var DB ArticleDB
	for rows.Next() {
		//rows.Scanの代わりにrows.StructScanを使う
		err := rows.StructScan(&DB)
		if err != nil {
			log.Fatal(err)
		}
		DB.Updatetime2 = StringTime(DB.Updatetime)
		resultStruct = append(resultStruct, DB)
	}
	query = "SELECT COUNT(status) FROM articles WHERE status = $1"
	rows, err = db.Queryx(query, "Alive")
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
	indexStruct := make([]CountPage, 0)
	div := (count % 5)
	if div == 0 {
		div = (count / 5)
	} else {
		div = (count / 5) + 1
	}
	if count <= 5 || div < index {
		// fmt.Println("5件以内or異常ページ感知")
		return indexStruct, div
	} else {
		// fmt.Println("for文開始")
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
		// fmt.Println("ソート前スライス=", slice)
		sort.Sort(sort.IntSlice(slice))
		// fmt.Println("ソート後スライス=", slice)
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
	rows, err := db.Queryx("SELECT * FROM articles")
	if err != nil {
		log.Fatal(err)
	}
	var DB ArticleDB
	for rows.Next() {
		//rows.Scanの代わりにrows.StructScanを使う
		err := rows.StructScan(&DB)
		if err != nil {
			log.Fatal(err)
		}
		DB.Updatetime2 = StringTime(DB.Updatetime)
		result = append(result, DB)
	}
	return result, err
}

func GetArticleOne(db *sqlx.DB, id string) ([]ArticleDB, error) {
	fmt.Println("GetArticles")
	result := make([]ArticleDB, 0)
	//SELECTを実行。db.Queryの代わりにdb.Queryxを使う。
	query := "SELECT * FROM articles WHERE id = $1"
	rows, err := db.Queryx(query, id)
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
		DB.Updatetime2 = StringTime(DB.Updatetime)
		result = append(result, DB)
	}
	return result, err
}

// NewArticlePosts = INSERT
func PostArticleNewPages(db *sqlx.DB, title string) int {
	fmt.Println("PostArticleNewPages")
	query := "INSERT INTO articles (title,status) VALUES($1,$2) RETURNING id"
	stmt, err := db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
	id := 0
	err = stmt.QueryRow(title, "Alive").Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func DeleteArticlePages(db *sqlx.DB, id int) error {
	fmt.Println("DeleteArticlePages")
	//SELECTを実行。db.Queryの代わりにdb.Queryxを使う。
	query := "UPDATE articles SET status = $1 WHERE id = $2"
	stmt, err := db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("Dead", id)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func DropArticlePages(db *sqlx.DB, id int, status string) error {
	fmt.Println("DropArticlePages")
	if status == "Dead" {
		query := "DELETE FROM articles WHERE id = $1"
		stmt, err := db.Prepare(query)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(id)
		if err != nil {
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
		return nil
	}
	//SELECTを実行。db.Queryの代わりにdb.Queryxを使う。
	query := "UPDATE articles SET status = $1 WHERE id = $2"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(status, id)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func UpdateArticleTimes(db *sqlx.DB, id string) error {
	fmt.Println("UpdateArticleTimes")
	query := "UPDATE articles SET updatetime = $1 WHERE id = $2"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(time.Now(), id)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func PostSearchPages(db *sqlx.DB, title string) ([]ArticleDB, error) {
	fmt.Println("PostSearchPages")
	keyword := "%" + title + "%"
	result := make([]ArticleDB, 0)
	//SELECTを実行。db.Queryの代わりにdb.Queryxを使う。
	query := "SELECT * FROM articles WHERE title LIKE $1"
	rows, err := db.Queryx(query, keyword)
	if err != nil {
		log.Fatal(err)
	}
	var DB ArticleDB
	for rows.Next() {
		//rows.Scanの代わりにrows.StructScanを使う
		err := rows.StructScan(&DB)
		if err != nil {
			log.Fatal(err)
		}
		DB.Updatetime2 = StringTime(DB.Updatetime)
		result = append(result, DB)
	}
	return result, err
}

func StringTime(t time.Time) string {
	return t.Format(L)
}
