package product

import (
	"got/core"
	"got/register"
	"syd/model"
	"syd/service/productservice"
)

func Register()
func init() {
	register.Component(Register, &ProductColorSizeTable{})
}

// ________________________________________________________________________________
// Product ColorSize Table
// version1 get product from db
// version2 generate table from passed parameters.
//

type ProductColorSizeTable struct {
	core.Component

	Id      int
	Product *model.Product
}

func (p *ProductColorSizeTable) Setup() {
	p.Product = productservice.GetProduct(p.Id)
}
