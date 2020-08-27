package util

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

func InitDB() *sqlx.DB {
	db, err := sqlx.Open("postgres", "user=postgres password=Asdfjkl; dbname=web_go sslmode=disable")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
