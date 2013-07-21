package product

import (
	"got/core"
	"got/debug"
	"got/register"
	"got/route"
	"syd/dal"
	"syd/model"
	"syd/service/productservice"
)

func init() {
	register.Component(Register, &ProductList{})
}

type ProductList struct {
	core.Component
	Products []*model.Product
	Source   string `query:"source"` // return here
}

// NOTE: event name is case sensitive. Kill this when add cache.
func (p *ProductList) Ondelete(productId int) (string, string) {
	debug.Log("Delete Product %d", productId)
	productservice.DeleteProduct(productId)
	// TODO make this default redirect.
	return route.RedirectDispatch(p.Source, "/product/list")
}

// display: total stocks
func (p *ProductList) NStock(productId int) (sum int) {
	stockmap := dal.ListProductStocks(productId)
	if stockmap != nil {
		for _, stock := range *stockmap {
			if stock > 0 {
				sum += stock
			}
		}
	}
	return
}
