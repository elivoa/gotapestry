package stat

// Deprecated, TODO chagne this into angularjs module.
import (
	"fmt"
	"syd/model"
	"syd/service"
	"time"

	"github.com/elivoa/got/builtin/services"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"github.com/elivoa/gxl"
)

type TrendStat struct {
	core.Component

	Period     int `query:"period"`
	CombineDay int `query:"combineday"`
	Yearonyear int `query:"yearonyear"`

	DailySalesData  model.ProductSales
	DailySalesData2 model.ProductSales

	Paylogs []*model.PayLog

	Days int `query:"days"` // paylog days

	// old things p
	Products []*model.Product
	Source   string `query:"source"` // return here
}

// display: total stocks
func (p *TrendStat) Setup() {
	// fmt.Println("===================================== ", p.Yearonyear)
	endtime := time.Now().AddDate(0, 0, 1).UTC().Truncate(time.Hour * 24)
	if salesdata, err :=
		service.Product.StatDailySalesData(0, p.Period, p.CombineDay, endtime); err != nil {
		return
	} else {
		p.DailySalesData = salesdata
	}

	// Year on year
	if p.Yearonyear > 0 {
		endtime = endtime.AddDate(-1, 0, 0)
		if salesdata, err :=
			service.Product.StatDailySalesData(0, p.Period, p.CombineDay, endtime); err != nil {
			return
		} else {
			p.DailySalesData2 = salesdata
		}
	}

	// 修改显示标签
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

	// get Paylogs

	// get time
	// TimeFrom, TimeTo := gxl.NatureTimeRange(0, 0, -1) // TODO NOT-UTC
	fmt.Println(">>>>>>>>>>>>>>>>>> Day is :", p.Days)
	if p.Days <= 0 {
		p.Days = 1
	}
	start, end := gxl.UntilStartOfTomorrowRangeUTC(p.Days)
	fmt.Println(">>>>>>>>>>>>>>>>>> ", start, end, p.Days)
	if data, err := service.Account.ListPaysByTime(start, end); err != nil {
		panic(err)
	} else {
		p.Paylogs = data
	}

	return
}

// 忽略 CombineNode 合并节点参数!
func (p *TrendStat) OnPeriod(days int) *exit.Exit {
	url := services.Link.GeneratePageUrlWithContextAndQueryParameters("stat/trend",
		map[string]interface{}{
			"period":     days,
			"combineday": 0, // 0的目的是点击period时间间隔的时候清除合并点。
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

func (p *TrendStat) OnYearonyear(yoy int) *exit.Exit {
	url := services.Link.GeneratePageUrlWithContextAndQueryParameters("stat/trend",
		map[string]interface{}{
			"period":     p.Period,
			"combineday": p.CombineDay,
			"yearonyear": yoy,
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

func (p TrendStat) YearonyearClass(yoy int) string {
	if yoy == p.Yearonyear {
		return "current"
	}
	return "-"
}

func (p *TrendStat) SumPays() float64 {
	if p.Paylogs == nil {
		return 0
	}
	var sum float64
	for _, paylog := range p.Paylogs {
		sum += paylog.Delta
	}
	return sum
}
