package dal

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"got/db"
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
		"product_id=? and property_name=?")
	defer db.CloseStmt(stmt)
	if db.Err(err) {
		return
	}

	rows, err := stmt.Query(propertyId, propertiesName)
	defer db.CloseRows(rows)
	if db.Err(err) {
		return
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
