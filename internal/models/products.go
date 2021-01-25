package models

type Product struct {
	ID          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Amount      int64  `json:"amount" db:"amount"`
	Price       int64  `json:"price" db:"price"`
	Barcode     string `json:"barcode" db:"barcode"`
}
