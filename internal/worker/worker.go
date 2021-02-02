package worker

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ENSLERMAN/warehouse-back/internal/utils"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

const (
	timer       = 3
	insertProds = `insert into warehouse.product_history 
					(dt, product_id, name, description, amount, price, barcode, is_delete) values 
					($1, $2, $3, $4, $5, $6, $7, $8);`
	insertDispatches = `insert into warehouse.dispatch_history 
					(dt, dispatch_id, dispatch_date, emp_id, emp_fio, status_id, status_name, customer_id, customer_fio) 
					values 
					($1, $2, $3, $4, $5, $6, $7, $8, $9);`
	insertShipments = `insert into warehouse.shipment_history 
					(dt, ship_id, supplier_id, supplier_fio, employee_id, employee_fio, date, product_barcode, product_amount) 
					values 
					($1, $2, $3, $4, $5, $6, $7, $8, $9);`
)

func UploadingHistory(clickCnn *sql.DB, pgCnn *sql.DB) {
	for {
		if err := uploadProducts(clickCnn, pgCnn); err != nil {
			logrus.Errorf("cannot upload product history: %v", err)
		}
		if err := uploadDispatches(clickCnn, pgCnn); err != nil {
			logrus.Errorf("cannot upload dispatches history: %v", err)
		}
		if err := uploadShipments(clickCnn, pgCnn); err != nil {
			logrus.Errorf("cannot upload shipments history: %v", err)
		}
		time.Sleep(time.Minute * timer)
	}
}

//nolint:funlen
func uploadProducts(clickCnn *sql.DB, pgCnn *sql.DB) error {
	type product struct {
		Date      time.Time
		ProductID int64  `db:"product_id"`
		Name      string `db:"name"`
		Desc      string `db:"description"`
		Amount    int64  `db:"amount"`
		Price     int64  `db:"price"`
		Barcode   string `db:"barcode"`
		IsDelete  bool   `db:"is_delete"`
	}

	pgRes, err := pgCnn.Query(`
		select 
			product_id, name, description, amount, price, barcode, is_delete
		from warehouse.product_history;`)
	if err != nil {
		logrus.Errorf("cannot get product history: %v", err)
		return err
	}
	products := make([]product, 0)
	for pgRes.Next() {
		pr := new(product)
		if err = pgRes.Scan(
			&pr.ProductID,
			&pr.Name,
			&pr.Desc,
			&pr.Amount,
			&pr.Price,
			&pr.Barcode,
			&pr.IsDelete,
		); err != nil {
			logrus.Errorf("cannot scan product: %v", err)
			continue
		}
		pr.Date = utils.GetTimeNowByMoscow()
		products = append(products, *pr)
	}

	if len(products) == 0 {
		return errors.New("products history is empty")
	}

	var (
		tx, _   = clickCnn.Begin()
		stmt, _ = tx.Prepare(insertProds)
	)

	for _, v := range products {
		if _, err = stmt.Exec(
			v.Date,
			v.ProductID,
			v.Name,
			v.Desc,
			v.Amount,
			v.Price,
			v.Barcode,
			strconv.FormatBool(v.IsDelete),
		); err != nil {
			logrus.Errorf("cannot upload to clickhouse %v, err: %v", v, err)
			return err
		}
	}

	if err = stmt.Close(); err != nil {
		logrus.Errorf("cannot close stmt, %v", err)
		return err
	}
	if err = tx.Commit(); err != nil {
		logrus.Errorf("Cannot make tx, %v", err)
		return err
	}

	if _, err = pgCnn.Exec(`truncate table warehouse.product_history;`); err != nil {
		logrus.Errorf("Cannot truncate table with err: %v", err)
		return err
	}

	return nil
}

//nolint:funlen
func uploadDispatches(clickCnn *sql.DB, pgCnn *sql.DB) error { //nolint:dupl
	type dispatch struct {
		DisID        int64   `json:"dispatch_id" db:"dispatch_id"`
		DispatchDate string  `json:"dispatch_date" db:"dispatch_date"`
		EmpID        *int64  `json:"emp_id" db:"emp_id"`
		EmpSurname   *string `json:"emp_surname" db:"emp_surname"`
		EmpName      *string `json:"emp_name" db:"emp_name"`
		EmpPat       *string `json:"emp_pat" db:"emp_pat"`
		EmpFIO       string  `db:"emp_fio"`
		StatusID     int64   `json:"status_id" db:"status_id"`
		StatusName   string  `json:"status_name" db:"status_name"`
		CusID        int64   `json:"cus_id" db:"cus_id"`
		CusSurname   string  `json:"cus_surname" db:"cus_surname"`
		CusName      string  `json:"cus_name" db:"cus_name"`
		CusPat       string  `json:"cus_pat" db:"cus_pat"`
		CusFIO       string  `db:"customer_fio"`
	}

	result, err := pgCnn.Query(`select * from warehouse.get_history_dispatches();`)
	if err != nil {
		logrus.Errorf("cannot get dispatches, %v", err)
		return err
	}
	dispatches := make([]dispatch, 0)
	for result.Next() {
		dis := new(dispatch)
		if err = result.Scan(
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
			logrus.Errorf("cannot get dispatches, %v", err)
			return err
		}
		dispatches = append(dispatches, *dis)
	}

	if len(dispatches) == 0 {
		return errors.New("dispatches history is empty")
	}

	var (
		tx, _   = clickCnn.Begin()
		stmt, _ = tx.Prepare(insertDispatches)
	)

	dt := utils.GetTimeNowByMoscow()
	for _, v := range dispatches {
		if v.EmpSurname != nil {
			v.EmpFIO = fmt.Sprintf("%s %s %s", *v.EmpSurname, *v.EmpName, *v.EmpPat)
		}
		v.CusFIO = fmt.Sprintf("%s %s %s", v.CusSurname, v.CusName, v.CusPat)
		if _, err = stmt.Exec(
			dt,
			v.DisID,
			v.DispatchDate,
			v.EmpID,
			v.EmpFIO,
			v.StatusID,
			v.StatusName,
			v.CusID,
			v.CusFIO,
		); err != nil {
			logrus.Errorf("cannot upload to clickhouse %v, err: %v", v, err)
			return err
		}
	}

	if err = stmt.Close(); err != nil {
		logrus.Errorf("cannot close stmt, %v", err)
		return err
	}
	if err = tx.Commit(); err != nil {
		logrus.Errorf("Cannot make tx, %v", err)
		return err
	}

	if _, err = pgCnn.Exec(`truncate table warehouse.dispatch_history;`); err != nil {
		logrus.Errorf("Cannot truncate table with err: %v", err)
		return err
	}

	return nil
}

//nolint:funlen
func uploadShipments(clickCnn *sql.DB, pgCnn *sql.DB) error { //nolint:dupl
	type shipment struct {
		ID              int64  `json:"id" db:"id"`
		SupplierID      int64  `json:"supplier_id" db:"supplier_id"`
		SupplierSurname string `json:"supplier_surname" db:"supplier_surname"`
		SupplierName    string `json:"supplier_name" db:"supplier_name"`
		SupplierPat     string `json:"supplier_pat" db:"supplier_pat"`
		SupplierFIO     string `db:"supplier_fio"`
		EmployeeID      int64  `json:"employee_id" db:"employee_id"`
		EmployeeSurname string `json:"employee_surname" db:"employee_surname"`
		EmployeeName    string `json:"employee_name" db:"employee_name"`
		EmployeePat     string `json:"employee_pat" db:"employee_pat"`
		EmployeeFIO     string `db:"employee_fio"`
		Date            string `json:"date" db:"date"`
		ProductBarcode  string `json:"product_barcode" db:"product_barcode"`
		ProductAmount   int    `json:"product_amount" db:"product_amount"`
	}

	result, err := pgCnn.Query(`select * from warehouse.get_shipments_history();`)
	if err != nil {
		logrus.Errorf("cannot get dispatches, %v", err)
		return err
	}
	ships := make([]shipment, 0)
	for result.Next() {
		s := new(shipment)
		err = result.Scan(&s.ID, &s.SupplierID, &s.SupplierSurname, &s.SupplierName, &s.SupplierPat, &s.EmployeeID,
			&s.EmployeeSurname, &s.EmployeeName, &s.EmployeePat, &s.Date, &s.ProductBarcode, &s.ProductAmount,
		)
		if err != nil {
			logrus.Errorf("cannot get shipments, %v", err)
			return err
		}
		ships = append(ships, *s)
	}

	if len(ships) == 0 {
		return errors.New("dispatches history is empty")
	}

	var (
		tx, _   = clickCnn.Begin()
		stmt, _ = tx.Prepare(insertShipments)
	)

	dt := utils.GetTimeNowByMoscow()
	for _, v := range ships {
		v.EmployeeFIO = fmt.Sprintf("%s %s %s", v.EmployeeSurname, v.EmployeeName, v.EmployeePat)
		v.SupplierFIO = fmt.Sprintf("%s %s %s", v.SupplierSurname, v.SupplierName, v.SupplierPat)
		if _, err = stmt.Exec(
			dt,
			v.ID,
			v.SupplierID,
			v.SupplierFIO,
			v.EmployeeID,
			v.EmployeeFIO,
			v.Date,
			v.ProductBarcode,
			v.ProductAmount,
		); err != nil {
			logrus.Errorf("cannot upload to clickhouse %v, err: %v", v, err)
			return err
		}
	}

	if err = stmt.Close(); err != nil {
		logrus.Errorf("cannot close stmt, %v", err)
		return err
	}
	if err = tx.Commit(); err != nil {
		logrus.Errorf("Cannot make tx, %v", err)
		return err
	}

	if _, err = pgCnn.Exec(`truncate table warehouse.shipment_history;`); err != nil {
		logrus.Errorf("Cannot truncate table with err: %v", err)
		return err
	}

	return nil
}
