package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func BindValidationError(ctx *gin.Context, err error, description string) {
	var errJSON struct {
		Code        int    `json:"code"`
		Err         string `json:"error"`
		Description string `json:"description"`
	}
	if err != nil {
		errJSON.Err = err.Error()
	}
	errJSON.Code = http.StatusBadRequest
	if description != "" {
		errJSON.Description = description
	}
	ctx.AbortWithStatusJSON(errJSON.Code, errJSON)
}

func BindServiceError(ctx *gin.Context, err error, description string) {
	var errJSON struct {
		Code        int    `json:"code"`
		Err         string `json:"error"`
		Description string `json:"description"`
	}
	if err != nil {
		errJSON.Err = err.Error()
	}
	errJSON.Code = http.StatusInternalServerError
	if description != "" {
		errJSON.Description = description
	}
	ctx.AbortWithStatusJSON(errJSON.Code, errJSON)
}

func BindDatabaseError(ctx *gin.Context, err error, description string) {
	var errJSON struct {
		Code        int    `json:"code"`
		Err         string `json:"error"`
		Description string `json:"description"`
	}
	if err != nil {
		errJSON.Err = err.Error()
	}
	errJSON.Code = http.StatusUnprocessableEntity
	if description != "" {
		errJSON.Description = description
	}
	ctx.AbortWithStatusJSON(errJSON.Code, errJSON)
}

func BindNoContent(ctx *gin.Context) {
	ctx.String(http.StatusNoContent, "")
}

func BindData(ctx *gin.Context, obj interface{}) {
	ctx.JSON(http.StatusOK, obj)
}

func BindUnauthorized(ctx *gin.Context, err error, description string) {
	var errJSON struct {
		Code        int    `json:"code"`
		Err         string `json:"error"`
		Description string `json:"description"`
	}
	if err != nil {
		errJSON.Err = err.Error()
	}
	errJSON.Code = http.StatusUnauthorized
	if description != "" {
		errJSON.Description = description
	}
	ctx.JSON(http.StatusUnauthorized, errJSON)
}
