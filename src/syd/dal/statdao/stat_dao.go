/*

changed to new db.

*/
package statdao

import (
	"github.com/elivoa/got/db"
	_ "github.com/go-sql-driver/mysql"
	"syd/model"
)

// TodayStat returns statistics of latest n days.
// TODO: return the second parameter as error
func TodayStat(n int) ([]*model.SumStat, error) {
	conn := db.Connectp() // Panic on connection error
	defer db.CloseConn(conn)
	sql := `
select DATEDIFF(create_time,NOW()) as 'date', 
  sum(1) as 'norder',
  sum(total_count) as 'nsold',
  sum(total_price) as '总价' ` +
		"from `order` where " + `
  DATEDIFF(create_time,NOW()) > ?
  and type <> 2
group by DATEDIFF(create_time,NOW())
order by DATEDIFF(create_time,NOW()) asc
`
	stmt, err := conn.Prepare(sql)
	if db.Err(err) {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(-n)
	if db.Err(err) {
		return nil, err
	}
	defer db.CloseRows(rows) // use db.CloseRows or rows.Close()? Is rows always nun-nil?

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
