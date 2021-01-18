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
			Date       string    `json:"date" validate:"required,gt=0"`
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
