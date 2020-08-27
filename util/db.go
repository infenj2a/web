package util

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func InitDB() *sqlx.DB {
	db, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
