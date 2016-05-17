package productdao

import (
	"database/sql"
	"github.com/elivoa/got/db"
	_ "github.com/go-sql-driver/mysql"
	"syd/model"
	"time"
)

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
