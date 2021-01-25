package handlers

import (
	"database/sql"
	"github.com/ENSLERMAN/warehouse-back/internal/models"
	"github.com/ENSLERMAN/warehouse-back/internal/utils"
	"github.com/gin-gonic/gin"
)

func GetProducts(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		result, err := db.Query("select * from warehouse.products")
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
