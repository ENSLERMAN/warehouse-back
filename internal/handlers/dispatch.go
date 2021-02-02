package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/ENSLERMAN/warehouse-back/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"time"
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
		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			loc = time.UTC
		}

		_, err = time.ParseInLocation(time.RFC3339, dis.DateCreate, loc)
		if err != nil {
			utils.BindValidationError(ctx, err, "Ошибка валидации даты")
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
			utils.BindDatabaseError(ctx, err, "cannot close dispatch")
			return
		}

		utils.BindNoContent(ctx)
	}
}

func GetDispatches(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		type dispatch struct {
			DisID        int64   `json:"dispatch_id" db:"dispatch_id"`
			EmpID        *int64  `json:"emp_id" db:"emp_id"`
			EmpSurname   *string `json:"emp_surname" db:"emp_surname"`
			EmpName      *string `json:"emp_name" db:"emp_name"`
			EmpPat       *string `json:"emp_pat" db:"emp_pat"`
			StatusID     int64   `json:"status_id" db:"status_id"`
			StatusName   string  `json:"status_name" db:"status_name"`
			DispatchDate string  `json:"dispatch_date" db:"dispatch_date"`
			CusID        int64   `json:"cus_id" db:"cus_id"`
			CusSurname   string  `json:"cus_surname" db:"cus_surname"`
			CusName      string  `json:"cus_name" db:"cus_name"`
			CusPat       string  `json:"cus_pat" db:"cus_pat"`
		}

		result, err := db.Query(`select * from warehouse.get_dispatches();`)
		if err != nil {
			utils.BindDatabaseError(ctx, err, "cannot get dispatches")
			return
		}
		dispatches := make([]dispatch, 0)
		for result.Next() {
			dis := new(dispatch)
			if err := result.Scan(
				&dis.DisID,
				&dis.EmpID,
				&dis.EmpSurname,
				&dis.EmpName,
				&dis.EmpPat,
				&dis.StatusID,
				&dis.StatusName,
				&dis.DispatchDate,
				&dis.CusID,
				&dis.CusSurname,
				&dis.CusName,
				&dis.CusPat,
			); err != nil {
				utils.BindDatabaseError(ctx, err, "cannot get dispatches")
				return
			}
			dispatches = append(dispatches, *dis)
		}
		utils.BindData(ctx, dispatches)
	}
}

func GetHistoryDispatches(db *sql.DB, clickDB *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		type dispatch struct {
			DisID        int64   `json:"dispatch_id" db:"dispatch_id"`
			DispatchDate string  `json:"dispatch_date" db:"dispatch_date"`
			EmpID        *int64  `json:"emp_id" db:"emp_id"`
			EmpSurname   *string `json:"emp_surname" db:"emp_surname"`
			EmpName      *string `json:"emp_name" db:"emp_name"`
			EmpPat       *string `json:"emp_pat" db:"emp_pat"`
			EmpFIO       *string `json:"emp_fio" db:"emp_fio"`
			StatusID     int64   `json:"status_id" db:"status_id"`
			StatusName   string  `json:"status_name" db:"status_name"`
			CusID        int64   `json:"cus_id" db:"cus_id"`
			CusSurname   string  `json:"cus_surname" db:"cus_surname"`
			CusName      string  `json:"cus_name" db:"cus_name"`
			CusPat       string  `json:"cus_pat" db:"cus_pat"`
			CusFIO       string  `json:"cus_fio" db:"cus_fio"`
		}

		result, err := db.Query(`select * from warehouse.get_history_dispatches();`)
		if err != nil {
			utils.BindDatabaseError(ctx, err, "cannot get dispatches")
			return
		}
		dispatches := make([]dispatch, 0)
		for result.Next() {
			dis := new(dispatch)
			if err := result.Scan(
				&dis.DisID,
				&dis.DispatchDate,
				&dis.EmpID,
				&dis.EmpSurname,
				&dis.EmpName,
				&dis.EmpPat,
				&dis.StatusID,
				&dis.StatusName,
				&dis.CusID,
				&dis.CusSurname,
				&dis.CusName,
				&dis.CusPat,
			); err != nil {
				utils.BindDatabaseError(ctx, err, "cannot get dispatches")
				return
			}
			dispatches = append(dispatches, *dis)
		}

		resultClick, err := clickDB.Query(`
				select
					dispatch_id, dispatch_date, emp_id, emp_fio, status_id, status_name, customer_id, customer_fio
				from warehouse.dispatch_history;
			`)
		if err != nil {
			utils.BindDatabaseError(ctx, err, "cannot get dispatches from clickhouse")
			return
		}
		for resultClick.Next() {
			dis := new(dispatch)
			if err := resultClick.Scan(
				&dis.DisID,
				&dis.DispatchDate,
				&dis.EmpID,
				&dis.EmpFIO,
				&dis.StatusID,
				&dis.StatusName,
				&dis.CusID,
				&dis.CusFIO,
			); err != nil {
				utils.BindDatabaseError(ctx, err, "cannot get dispatches from clickhouse")
				return
			}
			dispatches = append(dispatches, *dis)
		}

		utils.BindData(ctx, dispatches)
	}
}

func GetProductsInDispatch(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		idString := ctx.Query("dis_id")
		if idString == "" {
			utils.BindValidationError(ctx, errors.New("query param 'dis_id' is required"), "")
			return
		}
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			utils.BindServiceError(ctx, err, err.Error())
			return
		}

		type products struct {
			CusID       int64  `json:"cus_id" db:"customer_id"`
			ProdID      int64  `json:"product_id" db:"product_id"`
			ProdAmount  int64  `json:"product_amount" db:"product_amount"`
			ProdName    string `json:"product_name" db:"product_name"`
			ProdDes     string `json:"product_des" db:"product_description"`
			ProrBarcode string `json:"product_barcode" db:"product_barcode"`
		}

		prods := make([]products, 0)
		result, err := db.Query(`select * from warehouse.get_products_by_dispatch($1);`, id)
		if err != nil {
			utils.BindDatabaseError(ctx, err, "cannot get products by dispatch "+idString)
			return
		}
		for result.Next() {
			prod := new(products)
			err := result.Scan(&prod.CusID, &prod.ProdID, &prod.ProdAmount, &prod.ProdName, &prod.ProdDes, &prod.ProrBarcode)
			if err != nil {
				utils.BindDatabaseError(ctx, err, "cannot get products by dispatch "+idString)
				return
			}
			prods = append(prods, *prod)
		}
		if err = result.Err(); err != nil {
			utils.BindDatabaseError(ctx, err, "cannot get products by dispatch "+idString)
			return
		}

		if len(prods) == 0 {
			utils.BindNoContent(ctx)
			return
		}

		utils.BindData(ctx, prods)
	}
}

func RefuseDispatch(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		type Dispatch struct {
			Date  string
			EmpID int64 `json:"emp_id" validate:"required,gt=0"`
			CusID int64 `json:"cus_id" validate:"required,gt=0"`
			DisID int64 `json:"dis_id" validate:"required,gt=0"`
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
		_, err = db.Exec(`select * from warehouse.refuse_dispatch($1, $2, $3, $4);`,
			&dis.DisID,
			&dis.EmpID,
			&dis.CusID,
			&dis.Date,
		)
		if err != nil {
			utils.BindDatabaseError(ctx, err, "cannot refuse dispatch")
			return
		}

		utils.BindNoContent(ctx)
	}
}
