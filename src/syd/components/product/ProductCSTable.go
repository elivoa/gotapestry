package product

import (
	"github.com/elivoa/got/core"
	"syd/model"
	"syd/service"
)

// ________________________________________________________________________________
// Product ColorSize Table
// version1 get product from db
// version2 generate table from passed parameters.
//

type ProductColorSizeTable struct {
	core.Component
	Tid       string // client id
	ProductId int    // product id
	Product   *model.Product
}

func (p *ProductColorSizeTable) Setup() {
	if product, err := service.Product.GetFullProduct(p.ProductId); err != nil {
		panic(err)
	} else {
		p.Product = product
	}
}
