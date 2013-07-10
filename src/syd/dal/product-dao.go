package dal

import (
	"got/db"
	"log"
	"syd/model"
	"time"
)

/*
   syd :: person
*/
/* create new item in db */
func CreateProduct(product *model.Product) *model.Product {
	db.Connect()
	defer db.Close()

	if logdebug {
		log.Printf("[dal] Create product: %v", product)
	}

	stmt, err := db.DB.Prepare("insert into product(name, productId, brand, price, supplier, factoryPrice, stock, note, pictures, createtime, updatetime) values(?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	result, err := stmt.Exec(product.Name, product.ProductId, product.Brand, product.Price, product.Supplier, product.FactoryPrice, product.Stock, product.Note, product.Pictures, time.Now(), time.Now())
	if err != nil {
		panic(err.Error())
	}
	lastInsertId, err := result.LastInsertId()
	if err == nil || lastInsertId > 0 {
		product.Id = int(lastInsertId)
	}
	return product
}

/* update an existing item */
func UpdateProduct(product *model.Product) {
	db.Connect()
	defer db.Close()

	if logdebug {
		log.Printf("[dal] Edit product: %v", product)
	}

	stmt, err := db.DB.Prepare("update product set name=?, productId=?, brand=?, price=?, supplier=? , factoryPrice=?, stock=?, note=?, pictures=?, updatetime=? where id=?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	stmt.Exec(product.Name, product.ProductId, product.Brand, product.Price, product.Supplier, product.FactoryPrice, product.Stock, product.Note, product.Pictures, time.Now(), product.Id)
}

// Get product [updated new version]
func GetProduct(id int) (*model.Product, error) {
	if logdebug {
		log.Printf("[dal] Get Product with id %v", id)
	}

	conn, _ := db.Connect()
	defer conn.Close()

	stmt, err := conn.Prepare("select * from product where id = ?")
	defer stmt.Close()
	if db.Err(err) {
		return nil, err
	}

	row := stmt.QueryRow(id)
	p := new(model.Product)
	err = row.Scan(&p.Id, &p.Name, &p.ProductId, &p.Brand, &p.Price, &p.Supplier, &p.FactoryPrice, &p.Stock, &p.Note, &p.Pictures, &p.CreateTime, &p.UpdateTime)
	if db.Err(err) {
		return nil, err
	}
	return p, nil
}

/*
  List person with type
  TODO pager/range
*/
func ListProduct() *[]model.Product {
	if logdebug {
		log.Printf("[dal] List all product")
	}

	db.Connect()
	defer db.Close()

	stmt, err := db.DB.Prepare("select * from product")
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

	db.Connect()
	defer db.Close()

	stmt, err := db.DB.Prepare("delete from product where id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	stmt.Exec(id)
}
