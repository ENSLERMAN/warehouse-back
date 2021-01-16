package handlers

import (
	"database/sql"
	"github.com/ENSLERMAN/warehouse-back/internal/models"
	"github.com/ENSLERMAN/warehouse-back/internal/utils"
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

func UpdateRole(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var user struct {
			UserID   int64 `json:"user_id" binding:"required"`
			AccessID int64 `json:"access_id" binding:"required"`
		}
		err := ctx.ShouldBindJSON(&user)
		if err != nil {
			utils.BindValidationError(ctx, err, "body validation error")
			return
		}

		_, err = db.Exec(`call warehouse.changeuserrole($1, $2);`,
			&user.UserID,
			&user.AccessID,
		)
		if err != nil {
			utils.BindDatabaseError(ctx, err, "cannot change user role")
			return
		}

		utils.BindNoContent(ctx)
	}
}
