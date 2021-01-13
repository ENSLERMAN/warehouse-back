package service

import (
	"database/sql"
	"github.com/sirupsen/logrus"
)

func initDb() {
	_, err := sql.Open("postgres", "")
	if err != nil {
		logrus.Fatalf("cannot connect to database: %v", err)
	}
}