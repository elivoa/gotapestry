package dal

import (
	_ "github.com/go-sql-driver/mysql"
	"got/db"
	"log"
	"syd/model"
	"time"
)

/* Set customer private price */
func SetCustomerPrice(personId int, productId int, price float64) {
	db.Connect()
	defer db.Close()

	if logdebug {
		log.Printf("[dal] Set customer price for %v on %v, $%v", personId, productId, price)
	}

	// first get price xxx. TODO performance.
	customerPrice := GetCustomerPrice(personId, productId)
	if customerPrice == nil {

		// create
		stmt, err := db.DB.Prepare("insert into customer_special_price " +
			"(person_id, product_id, price, create_time, last_used_time) " +
			"values(?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()
		stmt.Exec(personId, productId, price, time.Now(), nil)

	} else {

		// update
		stmt, err := db.DB.Prepare("update customer_special_price set " +
			"person_id=?, product_id=?, price=?, last_used_time=? " +
			"where id = ?")
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()
		stmt.Exec(personId, productId, price, time.Now(), customerPrice.Id)
	}

}

//
// Get Customer Price, if nothing found, return nil.
//
func GetCustomerPrice(personId int, productId int) *model.CustomerPrice {
	prices := getCustomerPrice(personId, productId, 1)
	if prices != nil && len(*prices) == 1 {
		price := (*prices)[0]
		if price.Id > 0 {
			return &price
		} else {
			return nil
		}
	}
	return nil
}

func GetCustomerPriceHistory(personId int, productId int) *[]model.CustomerPrice {
	prices := getCustomerPrice(personId, productId, 1)
	if prices != nil {
		return prices
	}
	return nil
}

func getCustomerPrice(personId int, productId int, number int) *[]model.CustomerPrice {
	conn, _ := db.Connect()
	defer conn.Close()

	stmt, err := db.DB.Prepare("select * from customer_special_price " +
		"where person_id = ? and product_id = ? order by create_time DESC limit ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query(personId, productId, number)
	defer db.CloseRows(rows)
	if db.Err(err) {
		return nil
	}

	ps := []model.CustomerPrice{}
	for rows.Next() {
		p := new(model.CustomerPrice)
		rows.Scan(&p.Id, &p.PersonId, &p.ProductId, &p.Price, &p.CreateTime, &p.LastUsedTime)
		ps = append(ps, *p)
	}
	return &ps
}
