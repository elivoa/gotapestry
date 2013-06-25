package dal

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"got/db"
	"log"
	"strings"
	"syd/model"
	"time"
)

/*_______________________________________________________________________________
  Create new item in db
  TODO: Add transaction support.
*/
func CreateOrder(order *model.Order) (*model.Order, error) {
	if logdebug {
		log.Printf("[dal] Create Order: %v", order)
	}

	fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	fmt.Println(order.Details)

	// special 000. create order.Details
	if order.Details != nil && len(order.Details) > 0 {
		order.CalculateOrder()

		// insert into db
		_, err := createOrderDetail(&order.Details)
		if db.Err(err) {
			return nil, err
		}
	}

	// 1. create connection.
	conn, _ := db.Connect()
	defer conn.Close()

	// 2. prepare ql
	stmt, err := conn.Prepare("insert into `order` (id, track_number, status, delivery_method, customer_id, total_price, total_count, price_cut, note, create_time, update_time) values(?,?,?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	if db.Err(err) {
		return nil, err
	}

	// 3. execute
	res, err := stmt.Exec(
		order.Id, order.TrackNumber, order.Status, order.DeliveryMethod,
		order.CustomerId, order.TotalPrice, order.TotalCount, order.PriceCut,
		order.Note, time.Now(), time.Now(),
	)
	if db.Err(err) {
		return nil, err
	}

	// 4. assign id
	liid, err := res.LastInsertId()
	order.Id = int(liid)

	return order, nil
}

func createOrderDetail(orderDetails *[]*model.OrderDetail) (*[]*model.OrderDetail, error) {
	// 1. connect
	conn, _ := db.Connect()
	defer conn.Close()

	// 2. prepare ql
	stmt, err := conn.Prepare("insert into `order_detail` (id, order_track_number, product_id, quantity, unit, selling_price, note) values(?,?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	// 3. execute
	for _, detail := range *orderDetails {
		if detail == nil {
			continue
		}
		res, err := stmt.Exec(detail.Id, detail.OrderTrackNumber, detail.ProductId, detail.Quantity, detail.Unit, detail.SellingPrice, detail.Note)
		if err != nil {
			return nil, err
		}
		// 4. assign id
		liid, err := res.LastInsertId()
		detail.Id = int(liid)
	}
	return orderDetails, nil
}

/*_______________________________________________________________________________

  Get order
*/
func GetOrder(id int) (*model.Order, error) {
	if logdebug {
		log.Printf("[dal] Get Order with id %v", id)
	}

	conn, _ := db.Connect()
	defer conn.Close()

	stmt, err := db.DB.Prepare("select * from `order` where id = ?")

	// stmt, err := db.DB.Prepare("select id, track_number, status, delivery_method, customer_id, total_price, total_count, price_cut, note, create_time, update_time, close_time from `order` where id = ?")
	defer stmt.Close()
	if db.Err(err) {
		return nil, err
	}

	row := stmt.QueryRow(id)
	p := new(model.Order)
	err = row.Scan(
		&p.Id, &p.TrackNumber, &p.Status, &p.DeliveryMethod, &p.CustomerId,
		&p.TotalPrice, &p.TotalCount, &p.PriceCut, &p.Note,
		&p.CreateTime, &p.UpdateTime, &p.CloseTime,
	)
	if db.Err(err) {
		return nil, err
	}
	log.Printf("--------------------------------------")
	log.Printf("id is: %v, p is:  %v", id, p)
	log.Printf("--------------------------------------")
	// cascade
	details, err := getOrderDetails(p.TrackNumber)
	if db.Err(err) {
		return nil, err
	}
	p.Details = *details
	return p, nil
}

func getOrderDetails(trackNumber int64) (*[]*model.OrderDetail, error) {
	conn, _ := db.Connect()
	defer conn.Close()

	stmt, err := db.DB.Prepare("select * from `order_detail` where order_track_number=? order by id asc")
	defer stmt.Close()
	if db.Err(err) {
		return nil, err
	}

	rows, err := stmt.Query(trackNumber)
	defer rows.Close()
	if db.Err(err) {
		return nil, err
	}

	// big performance issue, maybe.
	orders := []*model.OrderDetail{}
	log.Println("---------================")
	log.Println(trackNumber)
	for rows.Next() {
		p := new(model.OrderDetail)
		rows.Scan(&p.Id, &p.OrderTrackNumber, &p.ProductId, &p.Quantity, &p.Unit, &p.SellingPrice, &p.Note)
		log.Println("---------")
		orders = append(orders, p)
	}
	return &orders, nil
}

/*_______________________________________________________________________________
  update an existing item
*/
func UpdateOrder(order *model.Order) (*model.Order, error) {
	if logdebug {
		log.Printf("[dal] Edit Order: %v", order)
	}

	// organize order details. delete all and then add all.
	err := deleteDetails(order.TrackNumber)
	if db.Err(err) {
		return nil, err
	}

	// insert into db
	order.CalculateOrder()
	_, err = createOrderDetail(&order.Details)
	if db.Err(err) {
		return nil, err
	}

	// update order
	conn, _ := db.Connect()
	defer conn.Close()

	stmt, err := db.DB.Prepare("update `order` set id=?, track_number=?, status=?, delivery_method=?, customer_id=?, total_price=?, total_count=?, price_cut=?, note=?, create_time=?, update_time=?, close_time=? where id=?")
	defer stmt.Close()
	if db.Err(err) {
		return nil, err
	}

	_, err = stmt.Exec(
		order.Id, order.TrackNumber, order.Status, order.DeliveryMethod,
		order.CustomerId, order.TotalPrice, order.TotalCount, order.PriceCut, order.Note,
		order.CreateTime, time.Now(), order.CloseTime,
		order.Id,
	)
	if db.Err(err) {
		return nil, err
	}

	return order, nil
}

func deleteDetails(trackNumber int64) error {
	conn, _ := db.Connect()
	defer conn.Close()

	stmt, err := db.DB.Prepare("delete from order_detail where order_track_number = ?")
	defer stmt.Close()
	if db.Err(err) {
		return err
	}

	_, err = stmt.Exec(trackNumber)
	if db.Err(err) {
		return err
	}
	return nil
}

/*_______________________________________________________________________________
  List person with type, default by status.
  Note: do not contains Details.
*/
func ListOrder(status string) *[]model.Order {
	// debug log
	if logdebug {
		log.Printf("[dal] List order with type:%v", status)
	}

	// header declare
	var err error

	// connection, // TODO need a connection pool?
	db.Connect()
	defer db.Close()

	// 1. query
	var queryString string
	if strings.ToLower(status) == "all" {
		queryString = "select * from `order`"
	} else {
		queryString = "select * from `order` where status = ?"
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
		rows, err = stmt.Query()
	} else {
		rows, err = stmt.Query(status)
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
