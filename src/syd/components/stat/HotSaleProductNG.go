package stat

import (
	"fmt"
	"got/core"
	"syd/model"
	"syd/service/productservice"
	"syd/service/statservice"
)

type HotSaleProduct2 struct {
	core.Component
	Days     int
	HotSales *statservice.HotSales
	products map[int]*model.Product
}

func (p *HotSaleProduct2) Setup() {
	p.products = make(map[int]*model.Product)
	p.HotSales = statservice.CalcHotSaleProducts(0, 0, -p.Days)
	for _, hsp := range p.HotSales.HSProduct {
		product := productservice.GetProduct(hsp.ProductId)
		if product != nil {
			p.products[hsp.ProductId] = product
		}
	}
}

func (p *HotSaleProduct2) ProductName(productId int) string {
	if p, ok := p.products[productId]; ok {
		return p.Name
	} else {
		return fmt.Sprint("-", productId, "-")
	}
}
