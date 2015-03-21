package statservice

import (
	"fmt"
	"sort"
	"syd/dal/orderdao"
	"syd/service/orderservice"
	"time"
)

type HotSales struct {
	HSProduct HotSaleProducts
}

type HotSaleProduct struct {
	ProductId int
	Sales     int
	Specs     map[string]int
}

type HotSaleProducts []*HotSaleProduct

func (p HotSaleProducts) Len() int           { return len(p) }
func (p HotSaleProducts) Less(i, j int) bool { return p[i].Sales > p[j].Sales }
func (p HotSaleProducts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// TODO make this meaningful.
// TODO 使用sql的方式查询统计。
func CalcHotSaleProducts(years, months, days int) *HotSales {
	orders, err := orderdao.ListOrderByTime(timeRangeWeek(years, months, days))
	if err != nil {
		return nil
	}

	pmap := map[int]*HotSaleProduct{}
	for _, o := range orders {

		// TODO very bad performance
		order, _ := orderservice.GetOrder(o.Id)
		if order.Details != nil {
			for _, d := range order.Details {
				cskey := fmt.Sprintf("%v_%v", d.Color, d.Size)
				hsp, ok := pmap[d.ProductId]
				if !ok {
					hsp = &HotSaleProduct{d.ProductId, 0, make(map[string]int)}
					hsp.Specs[cskey] = 0
					pmap[d.ProductId] = hsp
				}
				hsp.Specs[cskey] += d.Quantity
				hsp.Sales += d.Quantity
			}
		}
	}
	hs := &HotSales{} //HSProduct: []*HotSaleProduct{}
	for _, hsp := range pmap {
		hs.HSProduct = append(hs.HSProduct, hsp)
	}
	sort.Sort(hs.HSProduct)

	// for _, hsp := range hs.HSProduct {
	// 	fmt.Println(hsp)
	// }
	return hs
}

func timeRangeWeek(years, months, days int) (start, end time.Time) {
	end = time.Now().Truncate(time.Hour * 24)
	start = end.AddDate(years, months, days)
	return
}
