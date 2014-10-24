/*

changed to new db.

*/
package statdao

import (
	"database/sql"
	"github.com/elivoa/got/db"
	_ "github.com/go-sql-driver/mysql"
	"syd/model"
)

// TodayStat returns statistics of latest n days.
// TODO: return the second parameter as error
func TodayStat(n int) ([]*model.SumStat, error) {
	var conn *sql.DB
	var stmt *sql.Stmt
	var err error
	if conn, err = db.Connect(); err != nil {
		return nil, err
	}
	defer conn.Close()

	_sql := `
select DATEDIFF(create_time,NOW()) as 'date', 
  sum(1) as 'norder',
  sum(total_count) as 'nsold',
  sum(total_price) as '总价' ` +
		"from `order` " + `
where 
  DATEDIFF(create_time,NOW()) > ?
  and type in (?,?)
  and status in (?,?,?,?)
group by DATEDIFF(create_time,NOW())
order by DATEDIFF(create_time,NOW()) asc
`
	if stmt, err = conn.Prepare(_sql); err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		-n,
		model.Wholesale, model.ShippingInstead,
		"toprint", "todeliver", "delivering", "done", // "canceled",// model.ToPrint, model.ToDeliver, model.Delivering, model.Done,
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
