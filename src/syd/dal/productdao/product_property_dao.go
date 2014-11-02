package productdao

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/elivoa/got/db"
	_ "github.com/go-sql-driver/mysql"
	"syd/model"
)

// Set customer private price, panic if any error occurs.
func AddProductProperty(productId int, propertyName string, property string) {
	conn := db.Connectp()
	defer db.CloseConn(conn)

	sql := "insert into product_property " +
		"(product_id, property_name, value) values(?,?,?)"

	stmt, err := conn.Prepare(sql)
	if db.Err(err) {
		panic(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(productId, propertyName, property)
	if db.Err(err) {
		panic(err.Error())
	}
}

// Delete Product Property by product, property name and value
func DeleteProductProperty(productId int, propertyName string, property string) {
	conn := db.Connectp()
	defer db.CloseConn(conn)

	stmt, err := conn.Prepare("delete from product_property where " +
		"product_id = ? and property_name = ? and value = ? limit 1")
	if db.Err(err) {
		panic(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(productId, propertyName, property)
	if db.Err(err) {
		panic(err.Error())
	}
}

//
// Delete Product Property by product, property name and value
//
func DeleteOneProductProperty(productId int, propertyName string) {
	// fmt.Println("______________________________________________________________")
	// fmt.Println(productId)
	// fmt.Println(propertyName)
	if productId <= 0 {
		panic("Error when DeleteOneProductProperty: productId: " + string(productId))
	}

	conn := db.Connectp()
	defer db.CloseConn(conn)

	stmt, err := conn.Prepare("delete from product_property where " +
		"product_id = ? and property_name = ?")
	if db.Err(err) {
		panic(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(productId, propertyName)
	if db.Err(err) {
		panic(err.Error())
	}
}

// get properties by product and properties' name
func GetProductProperties(propertyId int, propertiesName string) (values []string) {
	conn, _ := db.Connect()
	defer db.CloseConn(conn) // should use db.CloseConn or conn.Close()?

	stmt, err := conn.Prepare("select value from product_property where " +
		"product_id=? and property_name=? order by id asc")
	defer db.CloseStmt(stmt)
	if db.Err(err) {
		panic(err.Error())
		// should here be empty return?
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
	return values
}

// TODO implement IsPropertyExist().
// TODO is this used?
func IsPropertyExist(productId string, propertyName string, propertyValue string) bool {
	return false
}

// delete product properties && create new properties.
// TODO add transaction here.
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

// fill color & sizes to product list.
func FillProductPropertiesByIdSet(models []*model.Product) error {
	if nil == models || len(models) == 0 {
		return nil
	}

	var conn *sql.DB
	var stmt *sql.Stmt
	var err error
	if conn, err = db.Connect(); err != nil {
		return err
	}
	defer conn.Close()

	var _sql bytes.Buffer        // sql buffer
	var params = []interface{}{} // params
	var index = map[int]*model.Product{}
	_sql.WriteString("select product_id, property_name, value from product_property where ")
	_sql.WriteString("product_id in (")
	for idx, m := range models {
		if idx > 0 {
			_sql.WriteRune(',')
		}
		_sql.WriteRune('?')
		params = append(params, m.Id)
		index[m.Id] = m
	}
	_sql.WriteRune(')')

	if stmt, err = conn.Prepare(_sql.String()); err != nil {
		return err
	}
	defer stmt.Close()

	// 3. execute
	rows, err := stmt.Query(params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// execute
	var productId int
	var propertyName string
	var value string

	for rows.Next() {
		err := rows.Scan(&productId, &propertyName, &value)
		if err != nil {
			return err
		}

		if product, ok := index[productId]; ok {
			if propertyName == "color" {
				if product.Colors == nil {
					product.Colors = []string{}
				}
				product.Colors = append(product.Colors, value)
			} else if propertyName == "size" {
				if product.Sizes == nil {
					product.Sizes = []string{}
				}
				product.Sizes = append(product.Sizes, value)
			}
		}
	}
	return nil
}

// func extractIdset(models []*model.Product) map[int]bool {
// 	var idarray = []int64{}
// 	if idset != nil {
// 		for id, _ := range idset {
// 			idarray = append(idarray, id)
// 		}
// 	}

// }
