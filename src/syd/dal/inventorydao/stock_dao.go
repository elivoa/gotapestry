package inventorydao

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/elivoa/got/db"
	"github.com/elivoa/got/debug"
	_ "github.com/go-sql-driver/mysql"
	"syd/model"
)

// ________________________________________________________________________________
// product color-size special values.
//
// NOTE: 1. only stock used. price is not used here.
//

/* Set special value of product color*size: stock and unit prices. */
func SetProductStock(productId int, color string, size string, stock int) {
	setProductCSValue(productId, color, size, "stock", stock, 0)
}

// TODO 没有引用。这个标的price字段是不是没有用啊。 所以注释掉。
// func SetProductPrice(productId int, color string, size string, price float64) {
// 	setProductCSValue(productId, color, size, "price", 0, price)
// }

// private 
func setProductCSValue(productId int, color string, size string,
	field string, stock int, price float64) {

	conn := db.Connectp()
	defer db.CloseConn(conn)

	_sql := fmt.Sprintf("insert into product_sku (product_id, color, size, %v) values (?,?,?,?) on duplicate key update %v = ?", field, field)

	stmt, err := conn.Prepare(_sql)
	defer db.CloseStmt(stmt) // the safe way to close.
	if db.Err(err) {
		panic(err.Error())
	}

	if field == "stock" {
		_, err := stmt.Exec(productId, color, size, stock, stock)
		if err != nil {
			debug.Error(err)
		}
	} else if field == "price" {
		_, err := stmt.Exec(productId, color, size, price, stock)
		if err != nil {
			debug.Error(err)
		}
	}
}

// update stock with delta.
func UpdateProductStockWithDelta(productId int64, color string, size string, stock int) error {
	var conn *sql.DB
	var stmt *sql.Stmt
	var err error
	if conn, err = db.Connect(); err != nil {
		return err
	}
	defer conn.Close()

	var _sql = fmt.Sprintf("insert into product_sku (product_id, color, size, price, stock, create_time) values (?,?,?,0,?, now() ) on duplicate key update stock = stock + ?")

	if stmt, err = conn.Prepare(_sql); err != nil {
		return err
	}
	defer stmt.Close()

	// 3. execute
	_, err = stmt.Exec(productId, color, size, stock, stock)
	if err != nil {
		return err
	}
	return nil
}

func ClearProductStock(productId int) error {
	conn := db.Connectp()
	defer db.CloseConn(conn)

	stmt, err := conn.Prepare("delete from product_sku where product_id = ?")
	if db.Err(err) {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(productId)
	if db.Err(err) {
		return err
	}
	return nil
}

/*_______________________________________________________________________________
  List Product Stocks
*/
func ListProductStocks(productId int) map[string]int {
	var err error
	conn := db.Connectp()
	defer db.CloseConn(conn)

	// 1. query
	var queryString = "select color,size,stock from `product_sku` where product_id = ?"

	// 2. prepare
	stmt, err := conn.Prepare(queryString)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	// 3. query
	rows, err := stmt.Query(productId)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// 4. process results.
	// big performance issue, maybe. who knows.
	var (
		color  string
		size   string
		stock  int
		stocks = map[string]int{}
	)

	for rows.Next() {
		rows.Scan(&color, &size, &stock)
		stocks[fmt.Sprintf("%v__%v", color, size)] = stock
	}
	return stocks
}

func GetProductStocks(productId int64, color, size string) (int, error) {
	var err error
	conn := db.Connectp()
	defer db.CloseConn(conn)

	// 1. query
	var queryString = "select `stock` from `product_sku` where product_id = ? and color=? and size=? "

	// 2. prepare
	stmt, err := conn.Prepare(queryString)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	// 3. query
	rows, err := stmt.Query(productId, color, size)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// 4. process results.
	var stock int = 0
	if rows.Next() {
		rows.Scan(&stock)
	}
	return stock, nil
}

/*_______________________________________________________________________________
  Fill order Lists.
*/
func filter(productId int) *map[string]int {
	var err error
	conn := db.Connectp()
	defer db.CloseConn(conn)

	// 1. query
	var queryString = "select color,size,stock from `product_sku` where product_id = ?"

	// 2. prepare
	stmt, err := conn.Prepare(queryString)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	// 3. query
	rows, err := stmt.Query(productId)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// 4. process results.
	// big performance issue, maybe. who knows.
	var (
		color  string
		size   string
		stock  int
		stocks = map[string]int{}
	)

	for rows.Next() {
		rows.Scan(&color, &size, &stock)
		stocks[fmt.Sprintf("%v__%v", color, size)] = stock
	}
	return &stocks
}

// ----------------------------------------------------------------------------------------------------
// fill product.Stocks and product.Stock
func FillProductStocksByIdSet(models []*model.Product) error {
	if nil == models || len(models) == 0 {
		return nil
	}

	var idset = map[int64]bool{}
	for _, m := range models {
		idset[int64(m.Id)] = true
	}
	if allstocks, err := GetAllStocksByIdSet(idset); err != nil {
		return err
	} else {
		if nil != allstocks {
			for _, m := range models {
				if stock, ok := allstocks[int64(m.Id)]; ok {
					m.Stocks = stock
					m.Stock = stock.Total()
				}
			}
		}
	}
	return nil
}

// ----------------------------------------------------------------------------------------------------
// fill stocks to inventory
func GetAllStocksByIdSet(idset map[int64]bool) (map[int64]model.Stocks, error) {
	if nil == idset || len(idset) == 0 {
		return nil, nil
	}
	var conn *sql.DB
	var stmt *sql.Stmt
	var err error
	if conn, err = db.Connect(); err != nil {
		return nil, err
	}
	defer conn.Close()

	var _sql bytes.Buffer        // sql buffer
	var params = []interface{}{} // params
	// var index = map[int]*model.Product{}
	_sql.WriteString("select id, product_id, color, size, stock from product_sku where ")
	_sql.WriteString("product_id in (")
	var idx = 0
	for id, b := range idset {
		if b {
			if idx > 0 {
				_sql.WriteRune(',')
			}
			_sql.WriteRune('?')
			params = append(params, id)
			// index[id] = m
			idx += 1
		}
	}
	_sql.WriteRune(')')

	if stmt, err = conn.Prepare(_sql.String()); err != nil {
		return nil, err
	}
	defer stmt.Close()

	// 3. execute
	rows, err := stmt.Query(params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// execute
	var (
		id        int
		productId int
		color     string
		size      string
		// price     float32
		stock int
	)

	returns := map[int64]model.Stocks{} // productId -> color -> size : stock
	for rows.Next() {
		err := rows.Scan(&id, &productId, &color, &size /*&price,*/, &stock)
		if err != nil {
			return nil, err
		}
		colors, ok := returns[int64(productId)]
		if !ok {
			colors = model.NewStocks()
			returns[int64(productId)] = colors
		}
		colors.Set(color, size, stock)
	}
	return returns, nil
}
