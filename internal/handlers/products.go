package handlers

import (
	"database/sql"
	"errors"
	"github.com/ENSLERMAN/warehouse-back/internal/models"
	"github.com/ENSLERMAN/warehouse-back/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strconv"
)

func GetProducts(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		result, err := db.Query("select id, name, description, amount, price, barcode from warehouse.products where is_delete = false;")
		if err != nil {
			utils.BindDatabaseError(ctx, err, "cannot get products")
			return
		}
		products := make([]models.Product, 0)
		for result.Next() {
			pr := new(models.Product)
			err := result.Scan(&pr.ID, &pr.Name, &pr.Description, &pr.Amount, &pr.Price, &pr.Barcode)
			if err != nil {
				utils.BindDatabaseError(ctx, err, "cannot get products")
				return
			}
			products = append(products, *pr)
		}

		if result.Err() != nil {
			utils.BindDatabaseError(ctx, result.Err(), "cannot get products")
			return
		}

		utils.BindData(ctx, products)
	}
}

func GetProductsByID(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		idString := ctx.Query("prod_id")
		if idString == "" {
			utils.BindValidationError(ctx, errors.New("query param 'prod_id' is required"), "")
			return
		}
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			utils.BindServiceError(ctx, err, err.Error())
			return
		}
		result := db.QueryRow("select id, name, description, amount, price, barcode from warehouse.products where id = $1 and is_delete = false;", id)
		if result.Err() != nil {
			utils.BindDatabaseError(ctx, result.Err(), "cannot get product by ID")
			return
		}
		pr := new(models.Product)
		err = result.Scan(&pr.ID, &pr.Name, &pr.Description, &pr.Amount, &pr.Price, &pr.Barcode)
		if err != nil {
			utils.BindDatabaseError(ctx, err, "cannot get product by ID")
			return
		}

		if result.Err() != nil {
			utils.BindDatabaseError(ctx, result.Err(), "cannot get product by ID")
			return
		}

		utils.BindData(ctx, pr)
	}
}

func UpdateProduct(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var prod struct {
			ID    int64  `json:"id" validate:"required,gt=0"`
			Name  string `json:"name" validate:"required,gt=0"`
			Desc  string `json:"desc" validate:"required,gt=0"`
			Price int64  `json:"price" validate:"required,gt=0"`
		}
		err := ctx.ShouldBindJSON(&prod)
		if err != nil {
			utils.BindValidationError(ctx, err, "body validation error")
			return
		}
		validate := validator.New()
		err = validate.Struct(prod)
		if err != nil {
			utils.BindValidationError(ctx, err, "body validation error")
			return
		}

		_, err = db.Exec(`select * from warehouse.edit_product($1, $2, $3, $4);`,
			&prod.ID,
			&prod.Name,
			&prod.Desc,
			&prod.Price,
		)
		if err != nil {
			utils.BindDatabaseError(ctx, err, "cannot change product")
			return
		}

		utils.BindNoContent(ctx)
	}
}

func DeleteProductsByID(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		idString := ctx.Query("prod_id")
		if idString == "" {
			utils.BindValidationError(ctx, errors.New("query param 'prod_id' is required"), "")
			return
		}
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			utils.BindServiceError(ctx, err, err.Error())
			return
		}
		_, err = db.Exec("select * from warehouse.delete_product($1);", id)
		if err != nil {
			utils.BindDatabaseError(ctx, err, "cannot delete product by ID")
			return
		}
		utils.BindNoContent(ctx)
	}
}
