package util

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func InitDB() *sqlx.DB {
	// db, err := sqlx.Open("postgres", "host=ec2-54-91-178-234.compute-1.amazonaws.com user=zsendmswafsfvc dbname=d304v44t1lic6h sslmode=require password=f8286a9a23f6af42df5aa2ff7d80e3432f8e073c3ccb454c1ada162d5e45831a")
	db, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	return db
}
