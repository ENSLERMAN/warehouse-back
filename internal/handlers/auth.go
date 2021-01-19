package handlers

import (
	"database/sql"
	"github.com/ENSLERMAN/warehouse-back/internal/models"
	"github.com/ENSLERMAN/warehouse-back/internal/utils"
	"github.com/gin-gonic/gin"
)

func Login(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var user struct {
			Login    string `json:"login" db:"login" binding:"required"`
			Password string `json:"password" db:"password" binding:"required"`
		}
		err := ctx.ShouldBindJSON(&user)
		if err != nil {
			utils.BindValidationError(ctx, err, "body validation error")
			return
		}

		result := db.QueryRow(`select login, password from warehouse.users where login = $1;`, user.Login)
		if err := result.Err(); err != nil {
			utils.BindUnauthorized(ctx, err, "")
			return
		}
		dbUser := &models.User{}
		if err = result.Scan(&dbUser.Login, &dbUser.Password); err != nil {
			utils.BindUnauthorized(ctx, err, "error in scan sql row")
			return
		}

		if utils.CheckPasswordHash(user.Password, dbUser.Password) {
			utils.BindUnauthorized(ctx, nil, "user or password is incorrect")
			return
		}

		result1 := db.QueryRow(`select * from warehouse.showinfobyme($1);`, &user.Login)
		if result1.Err() != nil {
			utils.BindDatabaseError(ctx, result1.Err(), "cannot get user data")
			return
		}

		var dbRes struct {
			UserID     int64  `json:"user_id" db:"user_id"`
			Surname    string `json:"surname" db:"surname"`
			Name       string `json:"name" db:"name"`
			Pat        string `json:"patronymic" db:"patronymic"`
			Login      string `json:"login" db:"patronymic"`
			AccessID   int64  `json:"access_id" db:"access_id"`
			AccessName string `json:"access_name" db:"access_name"`
		}

		if err := result1.Scan(
			&dbRes.UserID,
			&dbRes.Surname,
			&dbRes.Name,
			&dbRes.Pat,
			&dbRes.Login,
			&dbRes.AccessID,
			&dbRes.AccessName,
		); err != nil {
			utils.BindDatabaseError(ctx, err, "cannot get data from db")
			return
		}

		utils.BindData(ctx, &dbRes)
	}
}

func Register(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var user struct {
			Surname    string `json:"surname" db:"surname" binding:"required"`
			Name       string `json:"name" db:"name" binding:"required"`
			Patronymic string `json:"patronymic" db:"patronymic" binding:"required"`
			Login      string `json:"login" db:"login" binding:"required"`
			Password   string `json:"password" db:"password" binding:"required"`
		}
		err := ctx.ShouldBindJSON(&user)
		if err != nil {
			utils.BindValidationError(ctx, err, "body validation error")
			return
		}

		password, err := utils.HashPassword(user.Password)
		if err != nil {
			utils.BindServiceError(ctx, err, "cannot hash password")
		}
		_, err = db.Exec(`call warehouse.register_user($1, $2, $3, $4, $5, $6);`,
			&user.Surname,
			&user.Name,
			&user.Patronymic,
			&user.Login,
			password,
			2,
		)
		if err != nil {
			utils.BindDatabaseError(ctx, err, "")
			return
		}

		utils.BindNoContent(ctx)
	}
}