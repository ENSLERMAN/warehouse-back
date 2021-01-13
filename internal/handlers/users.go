package handlers

import (
	"database/sql"
	"fmt"
	"github.com/ENSLERMAN/warehouse-back/internal/models"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
)

func GetAllUsers(c *gin.Context) {
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

	url := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=ensler sslmode=disable", dbURL, dbUser, dbPassword)
	db, err := sql.Open("postgres", url)
	if err != nil {
		logrus.Fatalf("cannot connect to database: %v", err)
	}
	if err := db.Ping(); err != nil {
		logrus.Fatalf("cannot ping database: %v", err)
		return
	}

	result, err := db.Query("select * from Warehouse.Users")
	if err != nil {
		logrus.Error(err)
	}
	users := make([]*models.User, 0)
	for result.Next() {
		user := new(models.User)
		err := result.Scan(&user.ID, &user.Surname, &user.Name, &user.Patronymic, &user.Login, &user.Password, &user.Access)
		if err != nil {
			logrus.Error(err)
		}
		users = append(users, user)
	}
	if err = result.Err(); err != nil {
		logrus.Error(err)
	}
	c.JSON(200, &users)
	//c.String(200, "get all users")
}

func GetUserByID(c *gin.Context) {
	c.String(200, "get user by id")
}
