package dal

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"got/db"
	"log"
	"strings"
	"syd/model"
)

// ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func ListOrderByCustomer(personId int, status string) *[]model.Order {
	if logdebug {
		log.Printf("[dal] List order by person %v, with type:%v", personId, status)
	}

	// header declare
	var err error

	// connection, // TODO need a connection pool?
	db.Connect()
	defer db.Close()

	// 1. query`
	var queryString string
	if strings.ToLower(status) == "all" {
		queryString = "select * from `order` where customer_id = ?"
	} else {
		queryString = "select * from `order` where customer_id = ? and status = ?"
	}

	// 2. prepare
	stmt, err := db.DB.Prepare(queryString)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	// 3. query
	var rows *sql.Rows
	if strings.ToLower(status) == "all" {
		rows, err = stmt.Query(personId)
	} else {
		rows, err = stmt.Query(personId, status)
	}
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// 4. process results.
	// big performance issue, maybe. who knows.
	orders := []model.Order{}
	for rows.Next() {
		p := new(model.Order)
		rows.Scan(
			&p.Id, &p.TrackNumber, &p.Status, &p.DeliveryMethod, &p.CustomerId,
			&p.TotalPrice, &p.TotalCount, &p.PriceCut, &p.Note,
			&p.CreateTime, &p.UpdateTime, &p.CloseTime,
		)
		orders = append(orders, *p)
	}
	return &orders
}

// later
func DeleteOrder(id int) error {
	if logdebug {
		log.Printf("[dal] delete person %n", id)
	}

	// 1. TODO delete details
	// deleteDetails()

	conn, _ := db.Connect()
	defer conn.Close()

	stmt, err := db.DB.Prepare("delete from person where id = ?")
	defer stmt.Close()
	if db.Err(err) {
		return err
	}

	_, err = stmt.Exec(id)
	if db.Err(err) {
		return err
	}
	return nil
}
