package product

// Deprecated, TODO chagne this into angularjs module.
import (
	"github.com/elivoa/got/builtin/services"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"syd/model"
	"syd/service"
)

type ProductSalesChart struct {
	core.Component

	ProductId  int `query:"productid"` // 在t:a 中的parameters参数用来传递。
	Period     int `query:"period"`
	CombineDay int `query:"combineday"`

	DailySalesData model.ProductSales

	// old things p
	Products []*model.Product
	Source   string `query:"source"` // return here
}

// display: total stocks
func (p *ProductSalesChart) Setup() {

	if salesdata, err :=
		service.Product.StatDailySalesData(p.ProductId, p.Period, p.CombineDay); err != nil {
		return
	} else {
		p.DailySalesData = salesdata
	}

	if p.Period == 0 {
		p.Period = 30
	}
	if p.CombineDay == 0 {
		switch p.Period {
		case 7:
			p.CombineDay = 1
		case 30:
			p.CombineDay = 5
		case 90:
			p.CombineDay = 7
		case 365:
			p.CombineDay = 7
		default:
			p.CombineDay = 1
		}
	}
	return
}

// 忽略 CombineNode 合并节点参数!
func (p *ProductSalesChart) OnPeriod(days int) *exit.Exit {
	url := services.Link.GeneratePageUrlWithContextAndQueryParameters("product/detail",
		map[string]interface{}{
			"period":     days,
			"combineday": 0, // 目的是点击period时间间隔的时候清除合并点策略
		}, p.ProductId,
	)
	return exit.Redirect(url)
}

func (p *ProductSalesChart) OnCombineNode(combine_day int) *exit.Exit {
	url := services.Link.GeneratePageUrlWithContextAndQueryParameters("product/detail",
		map[string]interface{}{
			"period":     p.Period,
			"combineday": combine_day,
		}, p.ProductId,
	)
	return exit.Redirect(url)
}

func (p ProductSalesChart) PeriodLinkClass(period int) string {
	if period == p.Period {
		return "current"
	}
	return "-"
}

func (p ProductSalesChart) CombineNodeClass(combine_day int) string {
	if combine_day == p.CombineDay {
		return "current"
	}
	return "-"
}
