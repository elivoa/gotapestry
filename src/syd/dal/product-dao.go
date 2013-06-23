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

	stmt, err := db.DB.Prepare("insert into product(name, productId, brand, price, supplier, factoryPrice, stock, note, createtime, updatetime) values(?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	result, err := stmt.Exec(product.Name, product.ProductId, product.Brand, product.Price, product.Supplier, product.FactoryPrice, product.Stock, product.Note, time.Now(), time.Now())
	if err != nil {
		panic(err.Error())
	}
	lastInsertId, err := result.LastInsertId()
	if err == nil || lastInsertId > 0 {
		//	product.Id = lastInsertId
	}
	return product
	// log.Printf("Create Product Result is: %v\n", result)
	// log.Println(err)
}

/* update an existing item */
func UpdateProduct(product *model.Product) {
	db.Connect()
	defer db.Close()

	if logdebug {
		log.Printf("[dal] Edit product: %v", product)
	}

	stmt, err := db.DB.Prepare("update product set name=?, productId=?, brand=?, price=?, supplier=? , factoryPrice=?, stock=?, note=?, updatetime=? where id=?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	stmt.Exec(product.Name, product.ProductId, product.Brand, product.Price, product.Supplier, product.FactoryPrice, product.Stock, product.Note, time.Now(), product.Id)
}

func GetProduct(id int) *model.Product {
	if logdebug {
		log.Printf("[dal] Get Product with id %v", id)
	}

	db.Connect()
	defer db.Close()

	stmt, err := db.DB.Prepare("select * from product where id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	p := new(model.Product)
	row.Scan(&p.Id, &p.Name, &p.ProductId, &p.Brand, &p.Price, &p.Supplier, &p.FactoryPrice, &p.Stock, &p.Note, &p.CreateTime, &p.UpdateTime)
	return p
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
		rows.Scan(&p.Id, &p.Name, &p.ProductId, &p.Brand, &p.Price, &p.Supplier, &p.FactoryPrice, &p.Stock, &p.Note, &p.CreateTime, &p.UpdateTime)
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
