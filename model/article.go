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
	Userid      string    `db:"userid"`
	Title       string    `db:"title"`
	Status      string    `db:"status"`
	Updatetime  time.Time `db:"updatetime"`
	Updatetime2 string
}

type CountPage struct {
	Page int
}

const (
	// TOPページの取得件数
	getRecord = 5
	// 時間のフォーマット用
	Layout = "2006-01-02 15:04:05"
	// 新規スレッド書き込み時のロケーション
	location = "Asia/Tokyo"
)

func GetArticles(db *sqlx.DB, index int) ([]ArticleDB, int, []CountPage, error) {
	fmt.Println("GetArticles")
	page := 0
	if index != 1 {
		page += (index - 1) * getRecord
	}
	resultStruct := make([]ArticleDB, 0)
	//SELECTを実行。db.Queryの代わりにdb.Queryxを使う。
	query := "SELECT * FROM articles WHERE status = $1 order by updatetime desc LIMIT $2 OFFSET $3"
	rows, err := db.Queryx(query, "Alive", getRecord, page)
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
		// 取得した時間を整形
		DB.Updatetime2 = StringTime(DB.Updatetime)
		resultStruct = append(resultStruct, DB)
	}

	//ArticlesDBの総件数を取得
	query = "SELECT COUNT(status) FROM articles WHERE status = $1"
	rows, err = db.Queryx(query, "Alive")
	//count情報の取得
	resultCount := CheckCount(rows)
	//count情報よりページャーの作成
	indexStruct, div := exprCount(resultCount, index)
	if err != nil {
		log.Fatal(err)
	}
	return resultStruct, div, indexStruct, err
}

func CheckCount(rows *sqlx.Rows) (count int) {
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
	// TOPページは5件ずつ出力の設定なので
	// 総数が40件ならページ総数は8必要 総数が41件ならページ総数は9必要
	// つまり、総数が5で割り切れない場合は+1をする
	div := (count % 5)
	if div == 0 {
		div = (count / 5)
	} else {
		div = (count / 5) + 1
	}
	// 5件未満もしくは異常値が入力されていた場合はページャーは生成しない
	if count <= 5 || div < index {
		return indexStruct, div
	} else {
		//基点となる数字を初期値として、前後の数字を追加していく
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
		//最大5個の数字が入ったスライスをソート
		sort.Sort(sort.IntSlice(slice))
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
	query := "SELECT * FROM articles"
	rows, err := db.Queryx(query)
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
		// 時間の整形
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
		// 時間の整形
		DB.Updatetime2 = StringTime(DB.Updatetime)
		result = append(result, DB)
	}
	return result, err
}

// NewArticlePosts = INSERT
func PostArticleNewPages(db *sqlx.DB, userid, title string) int {
	fmt.Println("PostArticleNewPages")
	fmt.Println("userid=", userid)
	query := "INSERT INTO articles (userid,title,status) VALUES($1,$2,$3) RETURNING id"
	stmt, err := db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
	id := 0
	// 挿入したデータのID値を取得
	err = stmt.QueryRow(userid, title, "Alive").Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func DeleteArticlePages(db *sqlx.DB, id int) error {
	fmt.Println("DeleteArticlePages")
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
	defer stmt.Close()
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
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	query := "UPDATE articles SET updatetime = $1 WHERE id = $2"
	stmt, err := db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(time.Now().In(loc), id)
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
		// 時間の整形
		DB.Updatetime2 = StringTime(DB.Updatetime)
		result = append(result, DB)
	}
	return result, err
}

// 時間の整形用
func StringTime(t time.Time) string {
	return t.Format(Layout)
}
