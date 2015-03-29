package service

import (
	"database/sql"
	"fmt"
	"github.com/elivoa/got/db"
	"sort"
	"syd/model"
	"time"
)

type StatService struct{}

// TODO make this meaningful.
// TODO 使用sql的方式查询统计。
// params: days - 0 means today
// func CalcHotSaleProducts(years, months, days int) *HotSales {

// 	start, end := natureTimeRange(years, months, days)
// 	query := service.Order.EntityManager().Select().Where().
// 		// Or("status", "delivering", "done").
// 		Or("type", model.Wholesale, model.SubOrder). // 排除代发大订单，统计子订单即可。
// 		Range("create_time", start, end).            // time range.
// 		Limit(10000)                                 // max 1w orders.
// 	orders, err := service.Order.ListOrders(query, service.WITH_PRODUCT)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil
// 	}

// 	pmap := map[int]*HotSaleProduct{}
// 	for _, o := range orders {
// 		// TODO very bad performance
// 		order := o
// 		// order, _ := orderservice.GetOrder(o.Id)
// 		if order.Details != nil {
// 			for _, d := range order.Details {
// 				cskey := fmt.Sprintf("%v_%v", d.Color, d.Size)
// 				hsp, ok := pmap[d.ProductId]
// 				if !ok {
// 					hsp = &HotSaleProduct{d.ProductId, 0, make(map[string]int)}
// 					hsp.Specs[cskey] = 0
// 					pmap[d.ProductId] = hsp
// 				}
// 				hsp.Specs[cskey] += d.Quantity
// 				hsp.Sales += d.Quantity
// 			}
// 		}
// 	}
// 	hs := &HotSales{} //HSProduct: []*HotSaleProduct{}
// 	for _, hsp := range pmap {
// 		hs.HSProduct = append(hs.HSProduct, hsp)
// 	}
// 	sort.Sort(hs.HSProduct)

// 	// for _, hsp := range hs.HSProduct {
// 	// 	fmt.Println(hsp)
// 	// }
// 	return hs
// }

// DAO service
func (s *StatService) CalculateHotSaleProducts(years, months, days int) (*model.HotSales, error) {
	start, end := natureTimeRange(years, months, days)

	var conn *sql.DB
	var stmt *sql.Stmt
	var err error
	if conn, err = db.Connect(); err != nil {
		return nil, err
	}
	defer conn.Close()

	_sql := "select product_id,p.name,sum(quantity) from `order_detail` od " +
		"left join product p on od.product_id = p.Id " +
		"where od.order_track_number in (" +
		"  SELECT o.track_number FROM `order` o WHERE (`type`=0 or `type`=2) " +
		"  and (`create_time`>=? and `create_time`<?) " +
		") group by product_id order by sum(quantity) desc"

	// fmt.Println(_sql)
	// fmt.Println("start: ", start)
	// fmt.Println("end  : ", end)
	if stmt, err = conn.Prepare(_sql); err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// collect results.
	hs := &model.HotSales{} //HSProduct: []*HotSaleProduct{}

	// models := []*model.HotSaleProduct{}
	for rows.Next() {
		var (
			quantity    = new(sql.NullInt64)
			productId   = new(sql.NullInt64)
			productName = new(sql.NullString)
		)
		if err := rows.Scan(productId, productName, quantity); err != nil {
			return nil, err
		}
		m := new(model.HotSaleProduct)
		if productId.Valid {
			m.ProductId = productId.Int64
		}
		if productName.Valid {
			m.ProductName = productName.String
		}
		if quantity.Valid {
			m.Sales = (int)(quantity.Int64)
		}
		hs.HSProduct = append(hs.HSProduct, m)
	}

	sort.Sort(hs.HSProduct)
	return hs, nil
}

func (s *StatService) CalculateHotSaleProducts_with_specs(years, months, days int) (*model.HotSales, error) {
	start, end := natureTimeRange(years, months, days)

	var conn *sql.DB
	var stmt *sql.Stmt
	var err error
	if conn, err = db.Connect(); err != nil {
		return nil, err
	}
	defer conn.Close()

	_sql := "select product_id,p.name,color,size,sum(quantity) from `order_detail` od " +
		"left join product p on od.product_id = p.Id " +
		"where od.order_track_number in (" +
		"  SELECT o.track_number FROM `order` o WHERE (`type`=0 or `type`=2) " +
		"  and (`create_time`>=? and `create_time`<?) " +
		") group by product_id,color,size order by sum(quantity) desc"

	// fmt.Println(_sql)
	// fmt.Println("start: ", start)
	// fmt.Println("end  : ", end)
	if stmt, err = conn.Prepare(_sql); err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// collect results.
	pmap := map[int64]*model.HotSaleProduct{}
	for rows.Next() {
		var (
			quantity    = new(sql.NullInt64)
			color       = new(sql.NullString)
			size        = new(sql.NullString)
			productId   = new(sql.NullInt64)
			productName = new(sql.NullString)
		)
		if err := rows.Scan(productId, productName, color, size, quantity); err != nil {
			return nil, err
		}
		m := new(model.HotSaleProduct)
		if productId.Valid {
			m.ProductId = productId.Int64
		}
		if productName.Valid {
			m.ProductName = productName.String
		}
		if quantity.Valid {
			m.Sales = (int)(quantity.Int64)
		}

		// combine specs.
		if color.Valid && size.Valid {
			cskey := fmt.Sprintf("%v_%v", color.String, size.String)
			hsp, ok := pmap[m.ProductId]
			if !ok {
				hsp = &model.HotSaleProduct{
					ProductId:   m.ProductId,
					ProductName: m.ProductName,
					Sales:       0,
					Specs:       make(map[string]int),
				}
				hsp.Specs[cskey] = 0
				pmap[m.ProductId] = hsp
			}
			hsp.Specs[cskey] += m.Sales
			hsp.Sales += m.Sales
		}

	}

	hs := &model.HotSales{}
	for _, hsp := range pmap {
		hs.HSProduct = append(hs.HSProduct, hsp)
	}
	sort.Sort(hs.HSProduct)
	return hs, nil
}

// bad
func timeRangeWeek(years, months, days int) (start, end time.Time) {
	end = time.Now().Truncate(time.Hour * 24)
	start = end.AddDate(years, months, days)
	return
}

func natureTimeRange(years, months, days int) (start, end time.Time) {
	natureEnd := time.Now().AddDate(0, 0, 1).UTC().Truncate(time.Hour * 24)
	natureStart := natureEnd.AddDate(years, months, days-1)
	return natureStart, natureEnd
}
