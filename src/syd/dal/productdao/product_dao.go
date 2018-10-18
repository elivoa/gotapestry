// refactored
package productdao

import (
	"database/sql"
	"fmt"
	"log"
	"syd/base"
	"syd/base/product"
	"syd/model"
	"time"

	"github.com/elivoa/got/config"
	"github.com/elivoa/got/db"
	_ "github.com/go-sql-driver/mysql"
)

var logdebug = true
var core_fields = []string{
	"name", "productId", "status", "brand", "price", "supplier", "factoryPrice", "discountPercent",
	"stock", "shelfno", "capital", "note", "pictures", "ProducePeriod", "createtime",
}
var em = &db.Entity{
	Table:        "product",
	PK:           "id",
	Fields:       append(append([]string{"id"}, core_fields...), "updatetime"),
	CreateFields: core_fields,
	UpdateFields: core_fields,
}

func init() {
	db.RegisterEntity("product", em)
}

func EntityManager() *db.Entity {
	return em
}

//
// Universal one and list private
//

func _one(query *db.QueryParser) (*model.Product, error) {
	m := new(model.Product)
	err := query.Query(
		func(rows *sql.Rows) (bool, error) {
			return false, rows.Scan(
				&m.Id, &m.Name, &m.ProductId, &m.Status, &m.Brand, &m.Price, &m.Supplier, &m.FactoryPrice, &m.DiscountPercent,
				&m.Stock, &m.ShelfNo, &m.Capital, &m.Note, &m.Pictures, &m.ProducePeriod, &m.CreateTime, &m.UpdateTime,
			)
		},
	)
	if err != nil {
		return nil, err
	}
	if m.Id > 0 {
		return m, nil
	}
	return nil, nil
}

func _list(query *db.QueryParser) ([]*model.Product, error) {
	models := make([]*model.Product, 0)
	if err := query.Query(
		func(rows *sql.Rows) (bool, error) {
			m := &model.Product{}
			err := rows.Scan(
				&m.Id, &m.Name, &m.ProductId, &m.Status, &m.Brand, &m.Price, &m.Supplier, &m.FactoryPrice, &m.DiscountPercent,
				&m.Stock, &m.ShelfNo, &m.Capital, &m.Note, &m.Pictures, &m.ProducePeriod, &m.CreateTime, &m.UpdateTime,
			)
			models = append(models, m)
			return true, err
		},
	); err != nil {
		return nil, err
	}
	return models, nil
}

//
// Universal Get and List Public
//

func Get(id int) (*model.Product, error) {
	return _one(em.Select().Where(em.PK, id))
}

func List(parser *db.QueryParser) ([]*model.Product, error) {
	// var query *db.QueryParser
	parser.SetEntity(em) // set entity manager into query parser.
	parser.Reset()       // to prevent if parser is used before. TODO:Is this necessary?
	// append default behavore.
	parser.DefaultOrderBy("createtime", db.DESC)
	parser.DefaultLimit(0, config.LIST_PAGE_SIZE)
	parser.Select()
	return _list(parser)
}

// func ListAll() ([]*model.Product, error) {
// 	return _list(em.Select().Limit(500))
// }

// func ListByCapital(capital string) ([]*model.Product, error) {
// 	return _list(em.Select().Where("capital", capital))
// }

// ________________________________________________________________________________
// Create person
//
func Create(product *model.Product) (*model.Product, error) {
	res, err := em.Insert().Exec(
		product.Name, product.ProductId, product.Status, product.Brand, product.Price, product.Supplier,
		product.FactoryPrice, product.DiscountPercent, product.Stock, product.ShelfNo, product.Capital,
		product.Note, product.ProducePeriod, product.Pictures, product.CreateTime,
	)
	if err != nil {
		return nil, err
	}
	liid, err := res.LastInsertId()
	product.Id = int(liid)
	return product, nil
}

func UpdateProduct(product *model.Product) (int64, error) {
	if logdebug {
		log.Printf("[dal] Update Product: %v", product)
	}
	// update order
	res, err := em.Update().Exec(
		product.Name, product.ProductId, product.Status, product.Brand, product.Price, product.Supplier,
		product.FactoryPrice, product.DiscountPercent, product.Stock, product.ShelfNo, product.Capital,
		product.Note, product.Pictures, product.ProducePeriod, product.CreateTime,
		product.Id,
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// ________________________________________________________________________________
// Delete a product
//
func Delete(id int) (int64, error) {
	res, err := em.Delete().Exec(id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// update stock with delta.
func ChangeStatus(productId int, status product.Status) (affacted int64, err error) {
	var conn *sql.DB
	var stmt *sql.Stmt
	if conn, err = db.Connect(); err != nil {
		return
	}
	defer conn.Close()

	var _sql = fmt.Sprintf("update product p set p.status = ? where p.id = ? limit 1")
	if stmt, err = conn.Prepare(_sql); err != nil {
		return
	}
	defer stmt.Close()

	// 3. execute
	_, err = stmt.Exec(status, productId)
	if err != nil {
		return
	}
	return
}

func ListProductsByIdSet(ids ...int64) (map[int64]*model.Product, error) {
	if nil == ids || len(ids) == 0 {
		return nil, nil
	}
	var query *db.QueryParser
	parser := em.Select().Where()
	query = parser.InInt64("id", ids...).OrderBy("id", db.DESC)

	models, err := _list(query)
	if err != nil {
		panic(err)
	}

	var modelmap = map[int64]*model.Product{}
	for _, u := range models {
		modelmap[(int64)(u.Id)] = u
	}
	return modelmap, nil
}

//
// 统计产品日销量
//
func DailySalesData(
	productId int, startTime string, excludeYangYi bool, endday time.Time) (
	model.ProductSales, error) {

	// skip to all data
	if 0 == productId {
		return DailySalesData_alldata(startTime, true)
	}

	var conn *sql.DB
	var stmt *sql.Stmt
	var err error
	if conn, err = db.Connect(); err != nil {
		return nil, err
	}
	defer conn.Close()

	_sql := `
select
  DATE_FORMAT(o.create_time, '%Y-%m-%d'), sum(od.quantity)
from
  order_detail od
  left join ` + "`" + `order` + "`" + ` o on od.order_track_number = o.track_number
where
  od.product_id=?
  and o.type in (?,?)
  and o.status in (?,?,?,?)
  and o.create_time >= ?
  and o.create_time <= ?
  and od.product_id <> ?
group by
  DATE_FORMAT(o.create_time, '%Y-%m-%d')
order by
  DATE_FORMAT(o.create_time, '%Y-%m-%d') asc
`
	if stmt, err = conn.Prepare(_sql); err != nil {
		return nil, err
	}
	defer stmt.Close()

	// query
	var excluded_product_id = 0
	if excludeYangYi {
		excluded_product_id = base.STAT_EXCLUDED_PRODUCT
	}

	rows, err := stmt.Query(
		productId,
		model.Wholesale, model.SubOrder, // model.ShippingInstead, // 查子订单
		"toprint", "todeliver", "delivering", "done",
		startTime, endday, excluded_product_id,
	)
	if db.Err(err) {
		return nil, err
	}
	defer rows.Close() // db.CloseRows(rows) // use db.CloseRows or rows.Close()? Is rows always nun-nil?

	// the final result
	// ps := []*model.SalesNode{}
	ps := model.ProductSales{}

	for rows.Next() {
		p := new(model.SalesNode)
		rows.Scan(&p.Key, &p.Value)
		ps = append(ps, p)
	}
	return ps, nil
}

func DailySalesData_alldata(startTime string, excludeYangYi bool) (model.ProductSales, error) {

	var conn *sql.DB
	var stmt *sql.Stmt
	var err error
	if conn, err = db.Connect(); err != nil {
		return nil, err
	}
	defer conn.Close()

	_sql := `
select
  DATE_FORMAT(o.create_time, '%Y-%m-%d'), sum(od.quantity)
from
  order_detail od
  left join ` + "`" + `order` + "`" + ` o on od.order_track_number = o.track_number
where
  o.type in (?,?)
  and o.status in (?,?,?,?)
  and o.create_time >= ?
  and od.product_id <> ?
group by
  DATE_FORMAT(o.create_time, '%Y-%m-%d')
order by
  DATE_FORMAT(o.create_time, '%Y-%m-%d') asc
`
	if stmt, err = conn.Prepare(_sql); err != nil {
		return nil, err
	}
	defer stmt.Close()

	var excluded_product_id = 0
	if excludeYangYi {
		excluded_product_id = base.STAT_EXCLUDED_PRODUCT
	}

	rows, err := stmt.Query(
		model.Wholesale, model.SubOrder, // model.ShippingInstead, // 查子订单
		"toprint", "todeliver", "delivering", "done",
		startTime, excluded_product_id,
	)

	if db.Err(err) {
		return nil, err
	}
	defer rows.Close() // db.CloseRows(rows) // use db.CloseRows or rows.Close()? Is rows always nun-nil?

	ps := model.ProductSales{}

	for rows.Next() {
		p := new(model.SalesNode)
		rows.Scan(&p.Key, &p.Value)
		ps = append(ps, p)
	}
	return ps, nil
}

// --------------------------------------------------------------------------------
// Product's top buyer list. In product/detail page.
func ProductBestBuyerList(productId int) (model.BestBuyerList, error) {
	var conn *sql.DB
	var stmt *sql.Stmt
	var err error
	if conn, err = db.Connect(); err != nil {
		return nil, err
	}
	defer conn.Close()

	_sql := `
select -- d.product_id, p.Name,
	o.customer_id, pp.Name, sum(d.quantity), d.selling_price ` +
		"from `order` o " +
		` right join order_detail d on d.order_track_number = o.track_number
	left join product p on d.product_id = p.Id
	left join person pp on o.customer_id=pp.Id
where
	-- and o.customer_id = 305
	d.product_id = ?
  and o.type in (?,?)
  and o.status in (?,?,?,?)
	-- and o.create_time >= "2015-08-14" -- and o.create_time < "2015-03-23 23:55:55"
	-- and o.track_number=1501161519337773 -- debug
group BY d.product_id, o.customer_id
order by sum(d.quantity) desc
`

	if stmt, err = conn.Prepare(_sql); err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		productId,
		model.Wholesale, model.SubOrder, // model.ShippingInstead, // 查子订单
		"toprint", "todeliver", "delivering", "done",
	)

	if db.Err(err) {
		return nil, err
	}
	defer rows.Close() // db.CloseRows(rows) // use db.CloseRows or rows.Close()? Is rows always nun-nil?

	ps := model.BestBuyerList{}
	for rows.Next() {
		p := new(model.BestBuyerListItem)
		rows.Scan(&p.CustomerId, &p.CustomerName, &p.Quantity, &p.SalePrice)
		ps = append(ps, p)
	}
	return ps, nil
}
