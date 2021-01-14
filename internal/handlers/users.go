package handlers

import (
	"database/sql"
	"github.com/ENSLERMAN/warehouse-back/internal/models"
	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"
)

func GetAllUsers(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
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
	}
}

func GetUserByID(c *gin.Context) {
	c.String(200, "get user by id")
}
