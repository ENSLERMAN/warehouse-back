package models

var Product struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Amount      int64  `json:"amount" binding:"required"`
	Price       int64  `json:"price" binding:"required"`
	Barcode     string `json:"barcode" binding:"required"`
}