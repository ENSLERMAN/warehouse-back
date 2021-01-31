package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/ENSLERMAN/warehouse-back/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
)

func AddNewShipment(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		type Product struct {
			Name        string `json:"name" validate:"required,gt=0"`
			Description string `json:"description" validate:"required,gt=0"`
			Amount      int64  `json:"amount" validate:"required,gt=0"`
			Price       int64  `json:"price" validate:"required,gt=0"`
			Barcode     string `json:"barcode" validate:"required,gt=0"`
		}

		type Shipment struct {
			Date       string    `json:"date"`
			SupplierID int64     `json:"supplier_id" validate:"required,gt=0"`
			EmpID      int64     `json:"emp_id" validate:"required,gt=0"`
			Products   []Product `json:"products" validate:"required,gt=0,dive"`
		}
		var ship = new(Shipment)
		err := ctx.ShouldBindJSON(&ship)
		if err != nil {
			utils.BindValidationError(ctx, err, "body validation error")
			return
		}

		validate := validator.New()
		err = validate.Struct(ship)
		if err != nil {
			utils.BindValidationError(ctx, err, "body validation error")
			return
		}
		ship.Date = utils.GetNowByMoscow()
		shipmentJSON, err := jsoniter.Marshal(&ship)
		if err != nil {
			utils.BindValidationError(ctx, err, "body validation error")
			return
		}

		rawMsg := json.RawMessage(shipmentJSON)
		_, err = db.Exec(`select * from warehouse.add_new_shipment($1);`, rawMsg)
		if err != nil {
			utils.BindDatabaseError(ctx, err, "cannot make new shipment")
			return
		}

		utils.BindNoContent(ctx)
	}
}

func GetAllShipments(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		type shipment struct {
			Id              int64  `json:"id" db:"id"`
			SupplierId      int64  `json:"supplier_id" db:"supplier_id"`
			SupplierSurname string `json:"supplier_surname" db:"supplier_surname"`
			SupplierName    string `json:"supplier_name" db:"supplier_name"`
			SupplierPat     string `json:"supplier_pat" db:"supplier_pat"`
			EmployeeId      int64  `json:"employee_id" db:"employee_id"`
			EmployeeSurname string `json:"employee_surname" db:"employee_surname"`
			EmployeeName    string `json:"employee_name" db:"employee_name"`
			EmployeePat     string `json:"employee_pat" db:"employee_pat"`
			Date            string `json:"date" db:"date"`
			ProductBarcode  string `json:"product_barcode" db:"product_barcode"`
			ProductAmount   int    `json:"product_amount" db:"product_amount"`
		}

		result, err := db.Query("select * from warehouse.get_shipments();")
		if err != nil {
			utils.BindDatabaseError(ctx, err, "cannot get shipments")
			return
		}
		ships := make([]shipment, 0)
		for result.Next() {
			s := new(shipment)
			err := result.Scan(&s.Id, &s.SupplierId, &s.SupplierSurname, &s.SupplierName, &s.SupplierPat, &s.EmployeeId,
				&s.EmployeeSurname, &s.EmployeeName, &s.EmployeePat, &s.Date, &s.ProductBarcode, &s.ProductAmount,
			)
			if err != nil {
				utils.BindDatabaseError(ctx, err, "cannot get shipments")
				return
			}
			ships = append(ships, *s)
		}
		if err = result.Err(); err != nil {
			utils.BindDatabaseError(ctx, result.Err(), "cannot get shipments")
		}

		utils.BindData(ctx, ships)
	}
}

func GetShipmentsHistory(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		type shipment struct {
			Id              int64  `json:"id" db:"id"`
			SupplierId      int64  `json:"supplier_id" db:"supplier_id"`
			SupplierSurname string `json:"supplier_surname" db:"supplier_surname"`
			SupplierName    string `json:"supplier_name" db:"supplier_name"`
			SupplierPat     string `json:"supplier_pat" db:"supplier_pat"`
			EmployeeId      int64  `json:"employee_id" db:"employee_id"`
			EmployeeSurname string `json:"employee_surname" db:"employee_surname"`
			EmployeeName    string `json:"employee_name" db:"employee_name"`
			EmployeePat     string `json:"employee_pat" db:"employee_pat"`
			Date            string `json:"date" db:"date"`
			ProductBarcode  string `json:"product_barcode" db:"product_barcode"`
			ProductAmount   int    `json:"product_amount" db:"product_amount"`
		}

		result, err := db.Query("select * from warehouse.get_shipments_history();")
		if err != nil {
			utils.BindDatabaseError(ctx, err, "cannot get shipments")
			return
		}
		ships := make([]shipment, 0)
		for result.Next() {
			s := new(shipment)
			err := result.Scan(&s.Id, &s.SupplierId, &s.SupplierSurname, &s.SupplierName, &s.SupplierPat, &s.EmployeeId,
				&s.EmployeeSurname, &s.EmployeeName, &s.EmployeePat, &s.Date, &s.ProductBarcode, &s.ProductAmount,
			)
			if err != nil {
				utils.BindDatabaseError(ctx, err, "cannot get shipments")
				return
			}
			ships = append(ships, *s)
		}
		if err = result.Err(); err != nil {
			utils.BindDatabaseError(ctx, result.Err(), "cannot get shipments")
		}

		utils.BindData(ctx, ships)
	}
}
