package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/ENSLERMAN/warehouse-back/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
)

func AddNewDispatch(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		type Product struct {
			Amount  int64  `json:"amount" validate:"required,gt=0"`
			Barcode string `json:"barcode" validate:"required,gt=0"`
		}

		type Dispatch struct {
			Date       string    `json:"date"`
			DateCreate string    `json:"date_create" validate:"required,gt=0"`
			CustomerID int64     `json:"customer_id" validate:"required,gt=0"`
			Products   []Product `json:"products" validate:"required,gt=0,dive"`
		}
		var dis = new(Dispatch)
		err := ctx.ShouldBindJSON(&dis)
		if err != nil {
			utils.BindValidationError(ctx, err, "body validation error")
			return
		}

		validate := validator.New()
		err = validate.Struct(dis)
		if err != nil {
			utils.BindValidationError(ctx, err, "body validation error")
			return
		}
		dis.Date = utils.GetNowByMoscow()
		disJSON, err := jsoniter.Marshal(&dis)
		if err != nil {
			utils.BindValidationError(ctx, err, "body validation error")
			return
		}

		rawMsg := json.RawMessage(disJSON)
		_, err = db.Exec(`select * from warehouse.add_new_dispatch($1);`, rawMsg)
		if err != nil {
			utils.BindDatabaseError(ctx, err, "cannot make new dispatch")
			return
		}

		utils.BindNoContent(ctx)
	}
}

func CloseDispatch(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		type Product struct {
			Amount  int64  `json:"amount" validate:"required,gt=0"`
			Barcode string `json:"barcode" validate:"required,gt=0"`
		}

		type Dispatch struct {
			Date       string    `json:"date"`
			DispatchID int64     `json:"dispatch_id" validate:"required,gt=0"`
			EmpID      int64     `json:"emp_id" validate:"required,gt=0"`
			CustomerID int64     `json:"customer_id" validate:"required,gt=0"`
			Products   []Product `json:"products" validate:"required,gt=0,dive"`
		}
		var dis = new(Dispatch)
		err := ctx.ShouldBindJSON(&dis)
		if err != nil {
			utils.BindValidationError(ctx, err, "body validation error")
			return
		}

		validate := validator.New()
		err = validate.Struct(dis)
		if err != nil {
			utils.BindValidationError(ctx, err, "body validation error")
			return
		}
		dis.Date = utils.GetNowByMoscow()
		disJSON, err := jsoniter.Marshal(&dis)
		if err != nil {
			utils.BindValidationError(ctx, err, "body validation error")
			return
		}

		rawMsg := json.RawMessage(disJSON)
		_, err = db.Exec(`select * from warehouse.close_dispatch($1);`, rawMsg)
		if err != nil {
			utils.BindDatabaseError(ctx, err, "cannot make new dispatch")
			return
		}

		utils.BindNoContent(ctx)
	}
}
