package service

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"os"
)

func initDb() {
	dbUser, exist := os.LookupEnv("DB_LOGIN")
	if !exist {
		logrus.Fatalf("DB_LOGIN not exist in .env")
	}
	dbPassword, exist := os.LookupEnv("DB_PASSWORD")
	if !exist {
		logrus.Fatalf("DB_PASSWORD not exist in .env")
	}
	dbURL, exist := os.LookupEnv("DB_URL")
	if !exist {
		logrus.Fatalf("DB_URL not exist in .env")
	}
	url := "postgres://"+ dbUser + ":" + dbPassword + "@" + dbURL
	db, err := sql.Open("postgres", url)
	if err != nil {
		logrus.Fatalf("cannot connect to database: %v", err)
	}
	if err := db.Ping(); err != nil {
		logrus.Fatalf("cannot ping database: %v", err)
		return
	}
}