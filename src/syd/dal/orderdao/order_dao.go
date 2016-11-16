package orderdao

import (
	"database/sql"
	"fmt"
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/db"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"syd/model"
	"time"
)

var logdebug = true
var orderFields = []string{
	"track_number", "status", "type", "customer_id",
	"delivery_method", "delivery_tracking_number", "express_fee", "shipping_address",
	"total_price", "total_count", "price_cut", "Accumulated",
	"note", "parent_track_number",
	"create_time", "update_time", "close_time",
}
var em = &db.Entity{
	Table:        "order",
	PK:           "id",
	Fields:       append([]string{"id"}, orderFields...),
	CreateFields: orderFields,
	UpdateFields: orderFields,
}

func EntityManager() *db.Entity {
	return em
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

// Create new Order item in db. Only OrderService call this. Not including details.
// TODO: Add transaction support.
func CreateOrder(order *model.Order) error {
	if logdebug {
		log.Printf("[dal] Create Order: %v", order)
	}

	// 1. create connection.
	res, err := em.Insert().Exec(
		order.TrackNumber, order.Status, order.Type, order.CustomerId,
		order.DeliveryMethod, order.DeliveryTrackingNumber, order.ExpressFee, order.ShippingAddress,
		order.TotalPrice, order.TotalCount, order.PriceCut, order.Accumulated,
		order.Note, order.ParentTrackNumber, time.Now(), time.Now(), time.Now(),
	)
	fmt.Println("======================================")
	fmt.Println("order.CreateTimeis ", time.Now())
	fmt.Println("======================================")

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

	// // organize order details. delete all and then add all.
	// if _, err := DeleteDetails(order.TrackNumber); err != nil {
	// 	return 0, err
	// }

	// // special 000. create order.Details
	// if order.Details != nil && len(order.Details) > 0 {
	// 	// insert into db
	// 	if err := CreateOrderDetail(order.Details); err != nil {
	// 		return 0, err
	// 	}
	// }

	// update order
	res, err := em.Update().Exec(
		order.TrackNumber, order.Status, order.Type, order.CustomerId,
		order.DeliveryMethod, order.DeliveryTrackingNumber, order.ExpressFee, order.ShippingAddress,
		order.TotalPrice, order.TotalCount, order.PriceCut, order.Accumulated,
		order.Note, order.ParentTrackNumber, order.CreateTime, time.Now(), order.CloseTime,
		order.Id,
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func UpdateOrderStatus(trackNumber int64, status string) (int64, error) {
	if logdebug {
		log.Printf("[dal] Update Order %v's Status to %v", trackNumber, status)
	}
	now := time.Now()
	res, err := em.Update("status", "update_time", "close_time").Where(
		"track_number", trackNumber).Exec(status, now, now)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// TODO execute many / batch insert
func CreateOrderDetail(orderDetails []*model.OrderDetail) error {
	if nil != orderDetails {
		for _, detail := range orderDetails {
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
	}
	return nil
}

func BatchUpdateOrderDetail(orderDetails []*model.OrderDetail) error {
	if nil != orderDetails {
		for _, detail := range orderDetails {
			if detail == nil {
				continue
			}
			res, err := detailem.Update().Exec(
				detail.OrderTrackNumber, detail.ProductId, detail.Color, detail.Size,
				detail.Quantity, detail.SellingPrice, detail.Note,
				detail.Id,
			)
			if err != nil {
				return err
			}
			liid, err := res.LastInsertId()
			detail.Id = int(liid)
		}
	}
	return nil
}

func DeleteOrderDetails(orderDetails []*model.OrderDetail) error {
	if nil != orderDetails {
		for _, detail := range orderDetails {
			if detail == nil {
				continue
			}
			res, err := detailem.Delete().Where(detailem.PK, detail.Id).Exec()
			if err != nil {
				return err
			}
			liid, err := res.LastInsertId()
			detail.Id = int(liid)
		}
	}
	return nil
}

func DeleteDetailsByTrackNumber(trackNumber int64) (int64, error) {
	if res, err := detailem.Delete().Where("order_track_number", trackNumber).Exec(); err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
}

/*_______________________________________________________________________________
  update an existing item
*/

/*_______________________________________________________________________________
  Get order
*/
func GetOrder(field string, value interface{}) (*model.Order, error) {
	p := new(model.Order)
	if err := em.Select().Where(field, value).Query(
		func(rows *sql.Rows) (bool, error) {
			return false, rows.Scan(
				&p.Id, &p.TrackNumber, &p.Status, &p.Type, &p.CustomerId,
				&p.DeliveryMethod, &p.DeliveryTrackingNumber, &p.ExpressFee, &p.ShippingAddress,
				&p.TotalPrice, &p.TotalCount, &p.PriceCut, &p.Accumulated,
				&p.Note, &p.ParentTrackNumber,
				&p.CreateTime, &p.UpdateTime, &p.CloseTime,
			)
		},
	); err != nil {
		return nil, err
	}
	if p.Id > 0 {
		// cascade
		details, err := GetOrderDetails(p.TrackNumber)
		if err != nil {
			return nil, err
		}
		p.Details = details
		return p, nil
	}
	// return nil, errors.New("Order not found!") // return error if not exists.p
	return nil, nil // ignore when entity not exits.
}

/*_______________________________________________________________________________
  Get order
*/

// CountOrder returns number of orders that are top level orders(i.e. not include suborders.)
// func CountOrder(status string) (int, error) {
// 	if logdebug {
// 		log.Printf("[dal] Count Order with Status to %v", status)
// 	}

// 	parser := em.Count().Where()
// 	if status != "all" {
// 		parser.And("status", status)
// 	}
// 	count, err := parser.Or("type", model.Wholesale, model.ShippingInstead).QueryInt()
// 	if err != nil {
// 		return -1, err
// 	}
// 	return count, nil
// }

func CountOrderByCustomer(status string, personId int) (int, error) {
	if logdebug {
		log.Printf("[dal] Count Order with Status to %v", status)
	}

	parser := em.Count().Where()
	if status != "all" {
		parser.And("status", status)
	}
	count, err := parser.
		Or("type", model.Wholesale, model.ShippingInstead).
		And("customer_id", personId).
		QueryInt()
	if err != nil {
		return -1, err
	}
	return count, nil
}

// list interface
// TODO Order by create time;
func GetOrderDetails(trackNumber int64) ([]*model.OrderDetail, error) {
	orders := make([]*model.OrderDetail, 0)
	err := detailem.Select().Where("order_track_number", trackNumber).OrderBy("product_id", db.DESC).Query(
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
	// stmt, err := conn?.Prepare("select * from `order_detail` where order_track_number=? order by id asc")
}

// TODO transaction
func DeleteOrder(trackNumber int64) (int64, error) {
	aff, erro := DeleteDetailsByTrackNumber(trackNumber)
	if erro != nil {
		return aff, erro
	}

	res, err := em.Delete().Where("track_number", trackNumber).Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

//
// ---- Basic List orders -------------------------------------------
//
func ListOrders(parser *db.QueryParser) ([]*model.Order, error) {
	// var query *db.QueryParser
	parser.SetEntity(em) // set entity manager into query parser.
	parser.Reset()       // to prevent if parser is used before. TODO:Is this necessary?
	// append default behavore.
	parser.DefaultOrderBy("create_time", db.DESC)
	parser.DefaultLimit(0, config.LIST_PAGE_SIZE)
	parser.Select()
	return _listOrder(parser)
}

func _listOrder(query *db.QueryParser) ([]*model.Order, error) {
	orders := make([]*model.Order, 0)
	if err := query.Query(
		func(rows *sql.Rows) (bool, error) {
			p := new(model.Order)
			err := rows.Scan(
				&p.Id, &p.TrackNumber, &p.Status, &p.Type, &p.CustomerId,
				&p.DeliveryMethod, &p.DeliveryTrackingNumber, &p.ExpressFee, &p.ShippingAddress,
				&p.TotalPrice, &p.TotalCount, &p.PriceCut, &p.Accumulated,
				&p.Note, &p.ParentTrackNumber,
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

//
// --------  Special List Order    --------------------------------------------------------------------
//

// in most list order func, SubOrder is Excluded.

func ListOrder(status string) ([]*model.Order, error) {
	var query *db.QueryParser
	parser := em.Select().Where()
	if status != "all" {
		parser.And("status", status)
	}
	query = parser.
		Or("type", model.Wholesale, model.ShippingInstead).
		OrderBy("create_time", db.DESC)
	return _listOrder(query)
}

// not inlude sub orders.
func ListOrderPager(status string, limit int, n int) ([]*model.Order, error) {
	var query *db.QueryParser
	if status == "all" {
		query = em.Select().Where().
			Or("type", model.Wholesale, model.ShippingInstead).
			OrderBy("create_time", db.DESC).
			Limit(limit, n)
	} else {
		query = em.Select().
			Where("status", status).
			Or("type", model.Wholesale, model.ShippingInstead).
			OrderBy("create_time", db.DESC).
			Limit(limit, n)
	}
	return _listOrder(query)
}

func ListOrderByType(orderType model.OrderType, status string) ([]*model.Order, error) {
	var query *db.QueryParser
	if status == "all" {
		query = em.Select().Where().And("type", orderType).OrderBy("create_time", db.DESC)
	} else {
		query = em.Select().
			Where("status", status).And("type", orderType).OrderBy("create_time", db.DESC)
	}
	return _listOrder(query)
}

// directly change to limit version.
func ListOrderByCustomer(personId int, status string, limit int, n int) ([]*model.Order, error) {
	parser := em.Select()
	if status != "all" {
		parser.And("status", status)
	}
	var query *db.QueryParser
	query = parser.And("customer_id", personId).
		Or("type", model.Wholesale, model.ShippingInstead).
		OrderBy("create_time", db.DESC).
		Limit(limit, n)
	return _listOrder(query)
}

// TODO ...
func ListOrderByCustomerToday(personId int, status string) ([]*model.Order, error) {
	var query *db.QueryParser
	if status == "all" {
		query = em.Select().Where("customer_id", personId).
			Or("type", model.Wholesale, model.ShippingInstead).
			OrderBy("create_time", db.DESC)
	} else {
		query = em.Select().Where("customer_id", personId).
			And("status", status).
			Or("type", model.Wholesale, model.ShippingInstead).
			OrderBy("create_time", db.DESC)
	}
	return _listOrder(query)
}

func ListSubOrders(trackNumber int64) ([]*model.Order, error) {
	var query = em.Select().Where("parent_track_number", trackNumber).And("type", model.SubOrder)
	return _listOrder(query)
}

func DeliveringUnclosedOrdersByCustomer(personId int) ([]*model.Order, error) {
	query := em.Select().Where("customer_id", personId).And("status", "delivering").Or("type", model.Wholesale, model.ShippingInstead)
	return _listOrder(query)
}

// list as following.
func ListOrderByCustomer_Time(customerId int, start, end time.Time) ([]*model.Order, error) {
	var query *db.QueryParser
	query = em.Select().Where().And("customer_id", customerId).
		Or("status", "todeliver", "delivering", "done").
		// And("type", model.Wholesale).
		Range("create_time", start, end)
	return _listOrder(query)
}
