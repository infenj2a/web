package util

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func InitDB() *sqlx.DB {
	conf := DBConfig{
		User:     "postgres",
		Password: "Asdfjkl;",
		Host:     "postgres",
		Port:     "5432",
		Database: "web_go",
	}
	param := conf.User + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" +
		conf.Database + "?parseTime=true&loc=Asia%2FTokyo&charset=utf8mb4"
	db, err := sqlx.Open("postgres", param)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
