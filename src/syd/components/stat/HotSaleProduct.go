package stat

import (
	"github.com/elivoa/got/core"
	"syd/model"
	"syd/service"
)

type HotSaleProduct struct {
	core.Component
	Days     int
	HotSales *model.HotSales
	// products map[int]*model.Product
}

func (p *HotSaleProduct) Setup() {
	var err error
	if p.HotSales, err = service.Stat.CalculateHotSaleProducts_with_specs(0, 0, -p.Days+1); err != nil {
		panic(err)
	}
}
