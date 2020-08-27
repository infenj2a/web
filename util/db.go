package util

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type DBConfig struct {
	User     string
	Password string
	Connect  string
	Port     string
	Database string
}

func InitDB() *sqlx.DB {
	conf := DBConfig{
		User:     "postgres",
		Password: "Asdfjkl;",
		Connect:  "disabel",
		Port:     "5432",
		Database: "web_go",
	}
	param := "user=" + conf.User + "dbname=" + conf.Database + "password=" + conf.Password + "sslmode=" + conf.Connect
	db, err := sqlx.Open("postgres", param)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
