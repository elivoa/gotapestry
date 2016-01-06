/*

changed to new db.

*/
package statdao

import (
	"database/sql"
	"fmt"
	"github.com/elivoa/got/db"
	_ "github.com/go-sql-driver/mysql"
	"syd/base"
	"syd/model"
	"time"
)

// TodayStat returns statistics of latest n days.
// TODO: return the second parameter as error
func TodayStat(startTime time.Time, n int) ([]*model.SumStat, error) {
	var debug_print_time = true

	var conn *sql.DB
	var stmt *sql.Stmt
	var err error
	if conn, err = db.Connect(); err != nil {
		return nil, err
	}
	defer conn.Close()

	startTime = startTime.UTC().Truncate(time.Hour*24).AddDate(0, 0, 1)
	endTime := startTime.AddDate(0, 0, -n).Truncate(time.Hour * 24)
	if debug_print_time {
		fmt.Println("((((())))) ----  start time:", startTime)
		fmt.Println("((((())))) ----  end   time:", endTime)
	}
	_sql := `
select DATE_FORMAT(o.create_time, '%Y-%m-%d') as 'date', 
  count(distinct o.track_number) as 'norder',
  sum(od.quantity) as 'nsold',
  sum(od.quantity * od.selling_price) as '总价' ` +
		"from `order` o " + `
  right join order_detail od on o.track_number = od.order_track_number
where
  o.create_time<?
  and o.create_time >= ?
  and DATEDIFF(o.create_time,?) > ?
  and o.type in (?,?)
  and o.status in (?,?,?,?)
  and od.product_id<>?
group by DATEDIFF(o.create_time,?)
order by DATEDIFF(o.create_time,?) asc
`
	if stmt, err = conn.Prepare(_sql); err != nil {
		return nil, err
	}
	defer stmt.Close()

	// now := time.Now()
	rows, err := stmt.Query(
		startTime,
		endTime,
		startTime, -n,
		model.Wholesale, model.ShippingInstead,
		"toprint", "todeliver", "delivering", "done",
		base.TODAY_STAT_EXCLUDED_PRODUCT,
		startTime,
		startTime,
	)
	if db.Err(err) {
		return nil, err
	}
	defer rows.Close() // db.CloseRows(rows) // use db.CloseRows or rows.Close()? Is rows always nun-nil?

	// the final result
	ps := []*model.SumStat{}
	for rows.Next() {
		p := new(model.SumStat)
		rows.Scan(&p.Id, &p.NOrder, &p.NSold, &p.TotalPrice)

		// update average.
		p.AvgPrice = p.TotalPrice / float64(p.NSold)

		ps = append(ps, p)
	}
	return ps, nil
}
