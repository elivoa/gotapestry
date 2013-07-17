package orderdao

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"got/db"
	"log"
	"strings"
	"syd/model"
	"time"
)

var logdebug = true
var orderFields = []string{
	"track_number", "status", "delivery_method", "express_fee", "customer_id", "total_price",
	"total_count", "price_cut", "Accumulated", "note", "create_time", "update_time", "close_time",
}
var em = &db.Entity{
	Table:        "order",
	PK:           "id",
	Fields:       append([]string{"id"}, orderFields...),
	CreateFields: orderFields,
	UpdateFields: orderFields,
}

var orderDetailFields = []string{
	"order_track_number", "product_id", "color", "size", "quantity", "selling_price", "note",
}

var detailem = &db.Entity{
	Table:        "order_detail",
	PK:           "id",
	Fields:       append([]string{"id"}, orderDetailFields...),
	CreateFields: orderDetailFields,
	UpdateFields: orderDetailFields,
}

func init() {
	db.RegisterEntity("order", em)
	db.RegisterEntity("orderdetail", detailem)
}

func ListOrder(status string) ([]*model.Order, error) {
	orders := make([]*model.Order, 0)
	var query *db.QueryParser
	if status == "all" {
		query = em.Select()
	} else {
		query = em.Select().Where("status", status)
	}
	if err := query.Query(
		func(rows *sql.Rows) (bool, error) {
			p := new(model.Order)
			err := rows.Scan(
				&p.Id, &p.TrackNumber, &p.Status, &p.DeliveryMethod, &p.ExpressFee,
				&p.CustomerId, &p.TotalPrice, &p.TotalCount, &p.PriceCut, &p.Accumulated, &p.Note,
				&p.CreateTime, &p.UpdateTime, &p.CloseTime,
			)
			orders = append(orders, p)
			return true, err
		},
	); err != nil {
		return nil, err
	}
	return orders, nil
}

/*_______________________________________________________________________________
  Create new item in db
  TODO: Add transaction support.
*/
func CreateOrder(order *model.Order) error {
	if logdebug {
		log.Printf("[dal] Create Order: %v", order)
	}

	// special 000. create order.Details
	if order.Details != nil && len(order.Details) > 0 {
		order.CalculateOrder()

		// insert into db
		if err := createOrderDetail(order.Details); err != nil {
			return err
		}
	}

	// 1. create connection.
	res, err := em.Insert().Exec(
		order.TrackNumber, order.Status, order.DeliveryMethod, order.ExpressFee,
		order.CustomerId, order.TotalPrice, order.TotalCount, order.PriceCut,
		order.Accumulated, order.Note, time.Now(), time.Now(), time.Now(),
	)
	if err != nil {
		return err
	}
	liid, err := res.LastInsertId()
	order.Id = int(liid)
	return nil
}

func UpdateOrder(order *model.Order) (int64, error) {
	if logdebug {
		log.Printf("[dal] Update Order: %v", order)
	}

	// organize order details. delete all and then add all.
	if _, err := deleteDetails(order.TrackNumber); err != nil {
		return 0, err
	}

	// special 000. create order.Details
	if order.Details != nil && len(order.Details) > 0 {
		order.CalculateOrder()

		// insert into db
		if err := createOrderDetail(order.Details); err != nil {
			return 0, err
		}
	}

	// update order
	res, err := em.Update().Exec(
		order.TrackNumber, order.Status, order.DeliveryMethod, order.ExpressFee,
		order.CustomerId, order.TotalPrice, order.TotalCount, order.PriceCut,
		order.Accumulated, order.Note, order.CreateTime, time.Now(), order.CloseTime,
		order.Id,
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// TODO execute many / batch insert
func createOrderDetail(orderDetails []*model.OrderDetail) error {
	for _, detail := range orderDetails {
		// fmt.Printf(">>> detail: %v=n", detail)
		if detail == nil {
			continue
		}
		res, err := detailem.Insert().Exec(
			detail.OrderTrackNumber, detail.ProductId, detail.Color, detail.Size,
			detail.Quantity, detail.SellingPrice, detail.Note,
		)
		if err != nil {
			return err
		}
		liid, err := res.LastInsertId()
		detail.Id = int(liid)
	}
	return nil
}

/*_______________________________________________________________________________
  update an existing item
*/

/*_______________________________________________________________________________
  Get order
*/
func GetOrder(id int) (*model.Order, error) {
	p := new(model.Order)
	if err := em.Select().Where("id", id).Query(
		func(rows *sql.Rows) (bool, error) {
			return false, rows.Scan(
				&p.Id, &p.TrackNumber, &p.Status, &p.DeliveryMethod, &p.ExpressFee,
				&p.CustomerId, &p.TotalPrice, &p.TotalCount, &p.PriceCut,
				&p.Accumulated, &p.Note,
				&p.CreateTime, &p.UpdateTime, &p.CloseTime,
			)
		},
	); err != nil {
		return nil, err
	}
	if p.Id > 0 {
		// cascade
		details, err := getOrderDetails(p.TrackNumber)
		if db.Err(err) {
			return nil, err
		}
		p.Details = details
		return p, nil
	}
	return nil, errors.New("Order not found!")
}

// list interface
// TODO Order by id asc
func getOrderDetails(trackNumber int64) ([]*model.OrderDetail, error) {
	orders := make([]*model.OrderDetail, 0)
	err := detailem.Select().Where("order_track_number", trackNumber).Query(
		func(rows *sql.Rows) (bool, error) {
			p := new(model.OrderDetail)
			err := rows.Scan(
				&p.Id, &p.OrderTrackNumber, &p.ProductId, &p.Color, &p.Size, &p.Quantity,
				&p.SellingPrice, &p.Note,
			)
			orders = append(orders, p)
			return true, err
		},
	)
	if err != nil {
		return nil, err
	}
	return orders, nil
	// stmt, err := db.DB.Prepare("select * from `order_detail` where order_track_number=? order by id asc")
}

func deleteDetails(trackNumber int64) (int64, error) {
	if res, err := detailem.Delete().Where("order_track_number", trackNumber).Exec(); err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
}

// --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

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
