package models

type Dispatch struct {
	ID         int64  `json:"id" db:"id"`
	Date       string `json:"date" db:"date"`
	EmpID      int64  `json:"emp_id" db:"emp_id"`
	StatusID   int64  `json:"status_id" db:"status_id"`
	DateCreate string `json:"date_create" db:"date_create"`
	CustomerID string `json:"customer_id" db:"customer_id"`
}
