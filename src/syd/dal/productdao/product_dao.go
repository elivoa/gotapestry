package productdao

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"got/db"
	"syd/model"
	"time"
)

var logdebug = true
var em = &db.Entity{
	Table: "product",
	PK:    "id",
	Fields: []string{"id", "name", "productId", "brand", "price", "supplier", "factoryPrice",
		"stock", "shelfno", "note", "pictures", "createtime", "updatetime"},
	CreateFields: []string{"name", "productId", "brand", "price", "supplier", "factoryPrice",
		"stock", "shelfno", "note", "pictures", "createtime", "updatetime"}, // CreateFields
}

func init() {
	db.RegisterEntity("product", em)
}

// ________________________________________________________________________________
// Get product
//
func Get(id int) (*model.Product, error) {
	p := new(model.Product)
	err := em.Select().Where("id", id).QueryOne(
		func(row *sql.Row) error {
			return row.Scan(
				&p.Id, &p.Name, &p.ProductId, &p.Brand, &p.Price, &p.Supplier, &p.FactoryPrice,
				&p.Stock, &p.ShelfNo, &p.Note, &p.Pictures, &p.CreateTime, &p.UpdateTime,
			)
		},
	)
	if err != nil {
		return nil, err
	}
	if p.Id > 0 {
		return p, nil
	}
	return nil, errors.New("Product not found!")
}

// personType: customer, factory
func ListAll() ([]*model.Product, error) {
	products := make([]*model.Product, 0)
	err := em.Select().Query(
		func(rows *sql.Rows) (bool, error) {
			p := new(model.Product)
			err := rows.Scan(
				&p.Id, &p.Name, &p.ProductId, &p.Brand, &p.Price, &p.Supplier, &p.FactoryPrice,
				&p.Stock, &p.ShelfNo, &p.Note, &p.Pictures, &p.CreateTime, &p.UpdateTime,
			)
			products = append(products, p)
			return true, err
		},
	)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// ________________________________________________________________________________
// Create person
//
func Create(product *model.Product) (*model.Product, error) {
	res, err := em.Insert().Exec(
		product.Name, product.ProductId, product.Brand, product.Price, product.Supplier, product.FactoryPrice,
		product.Stock, product.ShelfNo, product.Note, product.Pictures, time.Now(), time.Now(),
	)
	if err != nil {
		return nil, err
	}
	liid, err := res.LastInsertId()
	product.Id = int(liid)
	return product, nil
}
