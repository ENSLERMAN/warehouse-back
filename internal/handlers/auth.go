package handlers

import (
	"database/sql"
	"github.com/ENSLERMAN/warehouse-back/internal/utils"
	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	var user struct {
		Login    string `json:"login" db:"login" binding:"required"`
		Password string `json:"password" db:"password" binding:"required"`
	}
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		utils.BindValidationError(ctx, err, "body validation error")
		return
	}
	utils.BindNoContent(ctx)
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
		_, err = db.Exec(`call register_user($1, $2, $3, $4, $5, $6);`,
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
