// refactored
package productdao

import (
	"database/sql"
	"errors"
	"github.com/elivoa/got/db"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"syd/model"
	"time"
)

var logdebug = true
var em = &db.Entity{
	Table: "product",
	PK:    "id",
	Fields: []string{"id", "name", "productId", "brand", "price", "supplier", "factoryPrice",
		"stock", "shelfno", "capital", "note", "pictures", "createtime", "updatetime"},
	CreateFields: []string{"name", "productId", "brand", "price", "supplier", "factoryPrice",
		"stock", "shelfno", "capital", "note", "pictures", "createtime", "updatetime"},
	UpdateFields: []string{"name", "productId", "brand", "price", "supplier", "factoryPrice",
		"stock", "shelfno", "capital", "note", "pictures", "createtime", "updatetime"},
}

func init() {
	db.RegisterEntity("product", em)
}

// ________________________________________________________________________________
// Get product
//
func Get(id int) (*model.Product, error) {
	p := new(model.Product)
	err := em.Select().Where("id", id).Query(
		func(row *sql.Rows) (bool, error) {
			return false, row.Scan(
				&p.Id, &p.Name, &p.ProductId, &p.Brand, &p.Price, &p.Supplier, &p.FactoryPrice,
				&p.Stock, &p.ShelfNo, &p.Capital, &p.Note, &p.Pictures, &p.CreateTime, &p.UpdateTime,
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

func ListAll() ([]*model.Product, error) {
	return _listProducts(em.Select())
}

func ListByCapital(capital string) ([]*model.Product, error) {
	return _listProducts(em.Select().Where("capital", capital))
}

func _listProducts(queryparser *db.QueryParser) ([]*model.Product, error) {
	products := make([]*model.Product, 0)
	err := queryparser.Query(
		func(rows *sql.Rows) (bool, error) {
			p := new(model.Product)
			err := rows.Scan(
				&p.Id, &p.Name, &p.ProductId, &p.Brand, &p.Price, &p.Supplier, &p.FactoryPrice,
				&p.Stock, &p.ShelfNo, &p.Capital, &p.Note, &p.Pictures, &p.CreateTime, &p.UpdateTime,
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
		product.Name, product.ProductId, product.Brand, product.Price, product.Supplier,
		product.FactoryPrice, product.Stock, product.ShelfNo, product.Capital,
		product.Note, product.Pictures, time.Now(), time.Now(),
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
		product.Name, product.ProductId, product.Brand, product.Price, product.Supplier,
		product.FactoryPrice, product.Stock, product.ShelfNo, product.Capital,
		product.Note, product.Pictures, time.Now(), time.Now(),
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
