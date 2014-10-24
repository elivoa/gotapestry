package dal

import (
	"database/sql"
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
	if prices, err := getCustomerPrice(personId, productId, 1); err != nil {
		panic(err.Error())
	} else if prices != nil && len(prices) == 1 {
		price := prices[0]
		if price.Id > 0 {
			return price
		} else {
			return nil
		}
	}
	return nil
}

func GetCustomerPriceHistory(personId int, productId int) []*model.CustomerPrice {
	if prices, err := getCustomerPrice(personId, productId, 1); err != nil {
		panic(err.Error())
	} else if prices != nil {
		return prices
	}
	return nil
}

func getCustomerPrice(personId int, productId int, number int) ([]*model.CustomerPrice, error) {
	var conn *sql.DB
	var stmt *sql.Stmt
	var err error
	if conn, err = db.Connect(); err != nil {
		return nil, err
	}
	defer conn.Close()

	_sql := "select * from customer_special_price where person_id = ? and product_id = ? order by create_time DESC limit ?"
	if stmt, err = conn.Prepare(_sql); err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(personId, productId, number)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	models := []*model.CustomerPrice{}
	for rows.Next() {
		m := new(model.CustomerPrice)
		var blackhole sql.NullInt64
		err := rows.Scan(&m.Id, &m.PersonId, &m.ProductId, &m.Price, &m.CreateTime, &blackhole /*&m.LastUsedTime */)
		if err != nil {
			panic(err)
		}
		models = append(models, m)
	}
	return models, nil
}
