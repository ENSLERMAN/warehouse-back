package handlers

import (
	"database/sql"
	"errors"
	"github.com/ENSLERMAN/warehouse-back/internal/models"
	"github.com/ENSLERMAN/warehouse-back/internal/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

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

func GetUsersByAccess(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		idString := ctx.Query("access_id")
		if idString == "" {
			utils.BindValidationError(ctx, errors.New("query param 'access_id' is required"), "")
			return
		}
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			utils.BindServiceError(ctx, err, err.Error())
			return
		}
		users := make([]models.User, 0)
		result, err := db.Query(`select id, surname, name, patronymic, login, access
			from warehouse.users where access = $1 and is_delete = false;`, id,
		)
		if err != nil {
			utils.BindDatabaseError(ctx, err, "cannot get suppliers")
			return
		}
		for result.Next() {
			user := new(models.User)
			err := result.Scan(&user.ID, &user.Surname, &user.Name, &user.Patronymic, &user.Login, &user.Access)
			if err != nil {
				utils.BindDatabaseError(ctx, err, "cannot get suppliers")
				return
			}
			users = append(users, *user)
		}
		if err = result.Err(); err != nil {
			utils.BindDatabaseError(ctx, err, "cannot get suppliers")
			return
		}

		if len(users) == 0 {
			utils.BindNoContent(ctx)
			return
		}

		utils.BindData(ctx, users)
	}
}
