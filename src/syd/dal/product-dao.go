// refactored.
package dal

import (
	"github.com/elivoa/got/db"
	"log"
	"syd/model"
	"time"
)

/*
   syd :: person
*/
/* create new item in db */
func CreateProduct(product *model.Product) *model.Product {
	conn := db.Connectp()
	defer db.CloseConn(conn)

	if logdebug {
		log.Printf("[dal] Create product: %v", product)
	}

	stmt, err := conn.Prepare("insert into product(name, productId, brand, price, supplier, factoryPrice, stock, note, pictures, createtime, updatetime) values(?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	result, err := stmt.Exec(product.Name, product.ProductId, product.Brand, product.Price, product.Supplier, product.FactoryPrice, product.Stock, product.Note, product.Pictures, time.Now(), time.Now())
	if err != nil {
		panic(err.Error())
	}

	// process result -- fill in ID
	lastInsertId, err := result.LastInsertId()
	if err == nil || lastInsertId > 0 {
		product.Id = int(lastInsertId)
	}
	return product
}

/* update an existing item */
func UpdateProduct(product *model.Product) {
	conn := db.Connectp()
	defer db.CloseConn(conn)

	if logdebug {
		log.Printf("[dal] Edit product: %v", product)
	}

	stmt, err := conn.Prepare("update product set name=?, productId=?, brand=?, price=?, supplier=? , factoryPrice=?, stock=?, shelfno=?,note=?, pictures=?, updatetime=? where id=?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.ProductId, product.Brand, product.Price, product.Supplier, product.FactoryPrice, product.Stock, product.ShelfNo, product.Note, product.Pictures, time.Now(), product.Id)
	if db.Err(err) {
		panic(err.Error())
	}
}

/*
  List person with type
  TODO pager/range
*/
func ListProduct() *[]model.Product {
	if logdebug {
		log.Printf("[dal] List all product")
	}

	conn := db.Connectp()
	defer db.CloseConn(conn)

	stmt, err := conn.Prepare("select * from product")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// big performance issue, maybe.
	products := []model.Product{}
	for rows.Next() {
		p := new(model.Product)
		rows.Scan(&p.Id, &p.Name, &p.ProductId, &p.Brand, &p.Price, &p.Supplier, &p.FactoryPrice, &p.Stock, &p.Note, &p.Pictures, &p.CreateTime, &p.UpdateTime)
		// TODO Peformance problem.
		products = append(products, *p)
	}
	return &products
}

// TODO error detect
func DeleteProduct(id int) {
	if logdebug {
		log.Printf("[dal] delete product %n", id)
	}

	conn := db.Connectp()
	defer db.CloseConn(conn)

	stmt, err := conn.Prepare("delete from product where id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if db.Err(err) {
		panic(err.Error())
	}
}
