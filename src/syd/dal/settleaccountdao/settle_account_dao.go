package settleaccountdao

import (
	"database/sql"
	"github.com/elivoa/got/db"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"syd/model"
	"time"
)

var logdebug = true

func init() {
	db.RegisterEntity("syd/settle_account", em)
}

var core_fields = []string{
	"factory_id", "goods_description", "from_time", "settle_time", "should_pay",
	"paid", "note", "operator",
}

var em = &db.Entity{
	Table:        "factory_settle_account",
	PK:           "id",
	Fields:       append([]string{"id"}, core_fields...),
	CreateFields: core_fields,
	UpdateFields: core_fields,
}

func EntityManager() *db.Entity {
	return em
}

// ________________________________________________________________________________
// Get product..... the following should be changed.
//
// func GetIntValue(name string) (int64, error) {
// 	if ccc, err := Get("name", name); err != nil {
// 		return 0, err
// 	} else {
// 		return strconv.ParseInt(ccc.Value, 10, 64)
// 	}
// }

// func GetStringValue(name string) (string, error) {
// 	if ccc, err := Get("name", name); err != nil {
// 		return "", err
// 	} else {
// 		return ccc.Value, nil
// 	}
// }

func GetById(id int64) (*model.FactorySettleAccount, error) {
	return Get(em.PK, id)
}

// 搞什么鬼
func Get(field string, value interface{}) (*model.FactorySettleAccount, error) {
	var query = em.Select().Where(field, value)
	if models, err := _list(query); err != nil {
		panic(err)
	} else {
		if nil != models && len(models) > 0 {
			if models[0] == nil {
				return nil, nil
			}
		}
		return models[0], nil
	}
}

func List(qp *db.QueryParser) ([]*model.FactorySettleAccount, error) {
	// parser := em.NewQueryParser().Select().Where("name", name)
	// parser.DefaultOrderBy("create_time", db.DESC)
	// parser.DefaultLimit(0, config.LIST_PAGE_SIZE)
	return _list(qp)
}

// the last part, read the list from rows
func _list(query *db.QueryParser) ([]*model.FactorySettleAccount, error) {
	models := make([]*model.FactorySettleAccount, 0)
	if err := query.Query(
		func(rows *sql.Rows) (bool, error) {
			m := &model.FactorySettleAccount{}
			err := rows.Scan(
				&m.Id, &m.FactoryId, &m.GoodsDescription, &m.FromTime, &m.SettleTime,
				&m.ShouldPay, &m.Paid, &m.Note, &m.Operator,
			)
			models = append(models, m)
			return true, err
		},
	); err != nil {
		return nil, err
	}
	return models, nil
}

// func GetByNameKey(name, key string) (*model.FactorySettleAccount, error) {
// 	var query = em.Select().Where("name", name).And("key", key)
// 	if models, err := _list(query); err != nil {
// 		panic(err)
// 	} else {
// 		if nil != models && len(models) > 0 {
// 			if models[0] != nil {
// 				return models[0], nil
// 			}
// 		}
// 		return nil, nil
// 	}
// }

// func Set(name string, key string, value interface{}, floatValue float64) error {
// 	var conn *sql.DB
// 	var stmt *sql.Stmt
// 	var err error
// 	if conn, err = db.Connect(); err != nil {
// 		return err
// 	}
// 	defer conn.Close()

// 	sql := "insert into const(`name`, `key`, `value`, doublevalue) values(?,?,?,?) on duplicate key update `key`=?, `value`=?, doublevalue=?"
// 	if stmt, err = conn.Prepare(sql); err != nil {
// 		return err
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(name, key, value, floatValue, key, value, floatValue)
// 	if err != nil {
// 		return err
// 	}
// 	return nil // res.RowsAffected()
// }

// "factory_id", "goods_description", "from_time", "settle_time", "should_pay",
// "paid", "note", "operator",

func Create(model *model.FactorySettleAccount) error {
	// 1. create connection.
	res, err := em.Insert().Exec(
		model.FactoryId, model.GoodsDescription, model.FromTime, model.SettleTime,
		model.ShouldPay, model.Paid, model.Note, model.OperatorId,
	)
	if err != nil {
		return err
	}
	liid, err := res.LastInsertId()
	model.Id = liid
	return nil
}

func Update(model *model.FactorySettleAccount) (int64, error) {
	if logdebug {
		log.Printf("[dal] Update FactorySettleAccount: %v", model)
	}

	// update order
	res, err := em.Update().Exec(
		model.FactoryId, model.GoodsDescription, model.FromTime, model.SettleTime,
		model.ShouldPay, model.Paid, model.Note, model.OperatorId,
		model.Id,
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func DeleteById(id int64) (int64, error) {
	res, err := em.Delete().Exec(id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// func Delete(name string, key string) (int64, error) {
// 	var conn *sql.DB
// 	var stmt *sql.Stmt
// 	var err error
// 	if conn, err = db.Connect(); err != nil {
// 		return 0, err
// 	}
// 	defer conn.Close()

// 	sql := "delete from const where name=? and key=?"
// 	if stmt, err = conn.Prepare(sql); err != nil {
// 		return 0, err
// 	}
// 	defer stmt.Close()

// 	// execute
// 	res, err := stmt.Exec(name, key)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return res.RowsAffected()
// }

// func DeleteByName(name string) (int64, error) {
// 	var conn *sql.DB
// 	var stmt *sql.Stmt
// 	var err error
// 	if conn, err = db.Connect(); err != nil {
// 		return 0, err
// 	}
// 	defer conn.Close()

// 	sql := "delete from const where name=?"
// 	if stmt, err = conn.Prepare(sql); err != nil {
// 		return 0, err
// 	}
// 	defer stmt.Close()

// 	// execute
// 	res, err := stmt.Exec(name)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return res.RowsAffected()
// }

// ----------------------------------------------------------------------------------------------------

func SettleAccount(startTime, endTime time.Time, factoryId int64) (*model.ProductSalesTable, error) {
	var conn *sql.DB
	var stmt *sql.Stmt
	var err error
	if conn, err = db.Connect(); err != nil {
		return nil, err
	}
	defer conn.Close()

	startTime = startTime.UTC().Truncate(time.Hour * 24)
	endTime = endTime.UTC().Truncate(time.Hour*24).AddDate(0, 0, 1)

	_sql := `
select product_id, sum(stock), send_time
from inventory i
where 
  send_time >= ?
  and send_time <= ?
  and provider_id = ?
group by DATE_FORMAT(send_time, '%Y-%m-%d'), product_id
`
	if stmt, err = conn.Prepare(_sql); err != nil {
		return nil, err
	}
	defer stmt.Close()

	// now := time.Now()
	rows, err := stmt.Query(
		startTime,
		endTime,
		factoryId,
		// model.Wholesale, model.SubOrder, // model.ShippingInstead, // 查子订单
		// "toprint", "todeliver", "delivering", "done",
		// base.STAT_EXCLUDED_PRODUCT,
		// startTime,
		// startTime,
	)
	if db.Err(err) {
		return nil, err
	}
	defer rows.Close() // db.CloseRows(rows) // use db.CloseRows or rows.Close()? Is rows always nun-nil?

	// the final result
	ps := model.NewProductSalesTable()
	var (
		productId int64
		stock     int
		send_time time.Time
	)

	for rows.Next() {
		rows.Scan(&productId, &stock, &send_time)
		ps.Set(send_time.Format("2006-01-02"), productId, stock)
	}
	return ps, nil
}
