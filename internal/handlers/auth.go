package handlers

import (
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

func Register(ctx *gin.Context) {
	var user struct {
		Surname    string `json:"surname" db:"surname" binding:"required"`
		Name       string `json:"name db:"name" binding:"required"`
		Patronymic string `json:"patronymic" db:"patronymic" binding:"required"`
		Login      string `json:"login" db:"login" binding:"required"`
		Password   string `json:"password" db:"password" binding:"required"`
		Access     uint64 `json:"access" db:"access" binding:"required"`
	}
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		utils.BindValidationError(ctx, err, "body validation error")
		return
	}

	_, err = utils.HashPassword(user.Password)
	if err != nil {
		utils.BindServiceError(ctx, err, "cannot hash password")
	}



	utils.BindNoContent(ctx)
}
