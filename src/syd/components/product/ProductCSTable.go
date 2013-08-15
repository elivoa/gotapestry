package product

import (
	"got/core"
	"syd/model"
	"syd/service/productservice"
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
	p.Product = productservice.GetProduct(p.ProductId)
}
