package product

import (
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"github.com/elivoa/gxl"
	"syd/model"
	"syd/service"
	"time"
)

/*
   Product List page
   -------------------------------------------------------------------------------
*/
type ProductSendNewProduct struct {
	core.Page

	// parameters
	RecentProductItems int64

	// table data
	Data *model.ProductSalesTable

	// others.

	// Products []*model.Product
	// Capital string `path-param:"1"`
	// ShowAll bool   `query:"showall"`
	Referer string `query:"referer"` // return here if non-empty

	TimeFrom time.Time `query:"from"`
	// TimeTo   time.Time `query:"to"`

	// temp daily lirun
	Date     string
	HotSales *model.Profits

	TotalCount           int
	TotalProfit          float64
	TotalProfitRite      float64
	PersonProfits        *model.PersonProfits
	NoPriceProductIdList string
}

func (p *ProductSendNewProduct) Setup() *exit.Exit {
	// // get items. TODO: Performance Issue.
	// recentProductItems, err := service.Const.Get2ndIntValue(SendNewProduct, RecentProductItems)
	// if err != nil {
	// 	panic(err)
	// }
	// p.RecentProductItems = recentProductItems

	// get data
	pst, err := service.SendNewProduct.GetSendNewProductTableData()
	if err != nil {
		panic(err)
	}
	p.Data = pst

	///	............
	// get daily product.
	p.Date = fmt.Sprint(p.TimeFrom)
	{
		var err error
		start, end := gxl.NatureTimeTodayRangeUTC(p.TimeFrom)

		if p.HotSales, err = service.Stat.CalculateHotSaleProducts_with_specs_specday(start, end); err != nil {
			panic(err)
		}

		// other tables.
		p.TotalCount = 0
		p.TotalProfit = 0
		p.TotalProfitRite = 0
		var count = 0
		var totalProfiteRiteSum float64 = 0
		for _, model := range p.HotSales.Profit {
			p.TotalProfit += model.Profit()
			p.TotalCount += model.Sales
			totalProfiteRiteSum += model.ProfitRate()
			fmt.Println("----------- ", p.TotalProfit, model.Profit(), model.Sales)
			if model.FactoryPrice > 0 { // 只计算有数据的，无数据的忽略
				count++
			} else {
				p.NoPriceProductIdList += fmt.Sprintf("%v,", model.ProductId)
			}
		}
		p.TotalProfitRite = totalProfiteRiteSum / float64(count)
		// fmt.Println(">>>>>>>>>>>>> ", totalProfiteRiteSum, count)

	}
	return nil
}

func (p *ProductSendNewProduct) ValidRow(profit *model.ProfitModel) string {
	if profit.FactoryPrice == 0 {
		return "color:red"
	}
	return ""
}

func (p *ProductSendNewProduct) TotalProfitRiteString() string {
	return fmt.Sprintf("%.1f %%", p.TotalProfitRite*100)
}
