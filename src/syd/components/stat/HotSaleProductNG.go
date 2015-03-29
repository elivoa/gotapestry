package stat

import (
	"github.com/elivoa/got/core"
	"syd/model"
	"syd/service"
)

type HotSaleProduct2 struct {
	core.Component
	Days     int
	HotSales *model.HotSales
}

func (p *HotSaleProduct2) Setup() {
	var err error
	if p.HotSales, err = service.Stat.CalculateHotSaleProducts_with_specs(0, 0, -p.Days+1); err != nil {
		panic(err)
	}
}
