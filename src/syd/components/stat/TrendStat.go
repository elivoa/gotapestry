package stat

// Deprecated, TODO chagne this into angularjs module.
import (
	"github.com/elivoa/got/builtin/services"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"syd/model"
	"syd/service"
)

type TrendStat struct {
	core.Component

	Period     int `query:"period"`
	CombineDay int `query:"combineday"`

	DailySalesData model.ProductSales

	// old things p
	Products []*model.Product
	Source   string `query:"source"` // return here
}

// display: total stocks
func (p *TrendStat) Setup() {

	if salesdata, err :=
		service.Product.StatDailySalesData(0, p.Period, p.CombineDay); err != nil {
		return
	} else {
		p.DailySalesData = salesdata
	}

	var no_argu = false
	if p.Period == 0 {
		no_argu = true
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
		case 1095:
			p.CombineDay = 30
		default:
			p.CombineDay = 1
		}
	}
	if no_argu {
		p.CombineDay = 1
	}
	return
}

// 忽略 CombineNode 合并节点参数!
func (p *TrendStat) OnPeriod(days int) *exit.Exit {
	url := services.Link.GeneratePageUrlWithContextAndQueryParameters("stat/trend",
		map[string]interface{}{
			"period":     days,
			"combineday": 0, // 目的是点击period时间间隔的时候清除合并点策略
		}, 0,
	)
	return exit.Redirect(url)
}

func (p *TrendStat) OnCombineNode(combine_day int) *exit.Exit {
	url := services.Link.GeneratePageUrlWithContextAndQueryParameters("stat/trend",
		map[string]interface{}{
			"period":     p.Period,
			"combineday": combine_day,
		}, 0,
	)
	return exit.Redirect(url)
}

func (p TrendStat) PeriodLinkClass(period int) string {
	if period == p.Period {
		return "current"
	}
	return "-"
}

func (p TrendStat) CombineNodeClass(combine_day int) string {
	if combine_day == p.CombineDay {
		return "current"
	}
	return "-"
}
