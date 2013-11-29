package dal

import (
	"github.com/elivoa/got/db"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"syd/model"
	"time"
)

/* Set customer private price */
func SetCustomerPrice(personId int, productId int, price float64) error {
	if logdebug {
		log.Printf("[dal] Set customer price for %v on %v, $%v", personId, productId, price)
	}

	conn := db.Connectp()
	defer db.CloseConn(conn)

	// first get price xxx. TODO performance.
	stmt, err := conn.Prepare("insert into customer_special_price " +
		"(person_id, product_id, price, create_time, last_use_time) " +
		"values(?,?,?,?,?)")
	if db.Err(err) {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(personId, productId, price, time.Now(), nil)
	if err != nil {
		return err
	}
	return nil
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
	conn := db.Connectp()
	defer db.CloseConn(conn)

	stmt, err := conn.Prepare("select * from customer_special_price " +
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
