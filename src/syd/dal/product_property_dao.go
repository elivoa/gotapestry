package dal

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"got/db"
	"got/debug"
)

/* Set customer private price */
func AddProductProperty(productId int, propertyName string, property string) {
	db.Connect()
	defer db.Close()

	// create
	stmt, err := db.DB.Prepare("insert into product_property " +
		"(product_id, property_name, value) values(?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	stmt.Exec(productId, propertyName, property)

	// // update
	// stmt, err := db.DB.Prepare("update customer_special_price set " +
	// 	"person_id=?, product_id=?, price=?, last_used_time=? " +
	// 	"where id = ?")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer stmt.Close()
	// stmt.Exec(personId, productId, price, time.Now(), customerPrice.Id)
}

//
// Delete Product Property by product, property name and value
//
func DeleteProductProperty(productId int, propertyName string, property string) {
	db.Connect()
	defer db.Close()

	stmt, err := db.DB.Prepare("delete from product_property where " +
		"product_id = ? and property_name = ? and value = ? limit 1")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	stmt.Exec(productId, propertyName, property)
}

//
// Delete Product Property by product, property name and value
//
func DeleteOneProductProperty(productId int, propertyName string) {
	fmt.Println("______________________________________________________________")
	fmt.Println(productId)
	fmt.Println(propertyName)
	if productId <= 0 {
		panic("Error when DeleteOneProductProperty: productId: " + string(productId))
	}

	db.Connect()
	defer db.Close()
	stmt, err := db.DB.Prepare("delete from product_property where " +
		"product_id = ? and property_name = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	stmt.Exec(productId, propertyName)
}

// get properties by product and properties' name
func GetProductProperties(propertyId int, propertiesName string) (values []string) {
	conn, _ := db.Connect()
	defer db.CloseConn(conn)

	stmt, err := conn.Prepare("select value from product_property where " +
		"product_id=? and property_name=? order by id asc")
	defer db.CloseStmt(stmt)
	if db.Err(err) {
		return
	}

	rows, err := stmt.Query(propertyId, propertiesName)
	defer db.CloseRows(rows)
	if db.Err(err) {
		return nil
	}

	// big performance issue, maybe.
	for rows.Next() {
		var propertyValue string
		rows.Scan(&propertyValue)
		values = append(values, propertyValue)
	}
	fmt.Println("????????????????????????????????????????????????")
	fmt.Println(values)
	return values
}

// TODO implement this.
func IsPropertyExist(productId string, propertyName string, propertyValue string) bool {
	return false
}

// delete product properties && create new properties.
func UpdateProductProperties(productId int, propertyName string, values ...string) {
	// properties := GetProperties(productId, propertyName)
	// TODO performance issue.
	DeleteOneProductProperty(productId, propertyName)
	if values != nil {
		fmt.Printf(">>>>>>>>>>>> properties: %v ; productId: %v\n", values, productId)
		for _, value := range values {
			AddProductProperty(productId, propertyName, value)
		}
	}
}

// ________________________________________________________________________________
// product color-size special values.
//
// NOTE: 1. only stock used. price is not used here.
//

/* Set special value of product color*size: stock and unit prices. */
func SetProductStock(productId int, color string, size string, stock int) {
	setProductCSValue(productId, color, size, "stock", stock, 0)
}

func SetProductPrice(productId int, color string, size string, price float64) {
	setProductCSValue(productId, color, size, "price", 0, price)
}

//   _________________
func setProductCSValue(productId int, color string, size string,
	field string, stock int, price float64) {

	db.Connect()
	defer db.Close()

	_sql := fmt.Sprintf("insert into product_cs_value (product_id, color, size, %v) values (?,?,?,?) on duplicate key update %v = ?", field, field)

	stmt, err := db.DB.Prepare(_sql)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

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

func ClearProductStock(productId int) error {
	conn, _ := db.Connect()
	defer conn.Close()

	stmt, err := db.DB.Prepare("delete from product_cs_value where product_id = ?")
	defer stmt.Close()
	if db.Err(err) {
		return err
	}

	_, err = stmt.Exec(productId)
	if db.Err(err) {
		return err
	}
	return nil
}

/*_______________________________________________________________________________
  List Product Stocks
*/
func ListProductStocks(productId int) *map[string]int {
	var err error
	db.Connect()
	defer db.Close()

	// 1. query
	var queryString = "select color,size,stock from `product_cs_value` where product_id = ?"

	// 2. prepare
	stmt, err := db.DB.Prepare(queryString)
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
	// fmt.Println(stocks)
	return &stocks
}
