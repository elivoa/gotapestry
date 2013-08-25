package statdao

import (
	_ "github.com/go-sql-driver/mysql"
	"got/db"
	"syd/model"
)

func TodayStat(n int) []*model.SumStat {
	conn, _ := db.Connect()
	defer conn.Close()
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

	stmt, err := db.DB.Prepare(sql)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query(-n)
	defer db.CloseRows(rows)
	if db.Err(err) {
		return nil
	}

	ps := []*model.SumStat{}
	for rows.Next() {
		p := new(model.SumStat)
		rows.Scan(&p.Id, &p.NOrder, &p.NSold, &p.TotalPrice)
		ps = append(ps, p)
	}
	return ps
}
