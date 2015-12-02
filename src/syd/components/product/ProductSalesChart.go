package product

// Deprecated, TODO chagne this into angularjs module.
import (
	"fmt"
	"github.com/elivoa/got/builtin/services"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"syd/model"
	"syd/service"
)

type ProductSalesChart struct {
	core.Component

	ProductId int `query:"productid"`
	Period    int

	DailySalesData model.ProductSales

	// old things p
	Products []*model.Product
	Source   string `query:"source"` // return here
}

// display: total stocks
func (p *ProductSalesChart) Setup() {

	if salesdata, err := service.Product.StatDailySalesData(p.ProductId, p.Period); err != nil {
		return
	} else {
		p.DailySalesData = salesdata
	}
	return
}

func (p *ProductSalesChart) OnPeriod(days int) *exit.Exit {
	url := services.Link.GeneratePageUrlWithContextAndQueryParameters("product/detail",
		map[string]interface{}{"period": days}, p.ProductId,
	)

	fmt.Println("\n-================= go to days ", days, p.ProductId)
	return exit.Redirect(url)
}

func (p ProductSalesChart) PeriodLinkClass(period int) string {
	if period == p.Period {
		return "current"
	}
	return "-"
}
