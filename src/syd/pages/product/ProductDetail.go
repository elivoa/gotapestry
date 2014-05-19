package product

import (
	"github.com/elivoa/gxl"
	"got/core"
	"syd/model"
	"syd/service/personservice"
	"syd/service/productservice"
)

// ________________________________________________________________________________
// Product Details

type ProductDetail struct {
	core.Page
	Id      *gxl.Int `path-param:"1"`
	Product *model.Product
}

func (p *ProductDetail) Setup() {
	p.Product = productservice.GetProduct(p.Id.Int)
}

func (p *ProductDetail) Pictures() []string {
	return productservice.ProductPictrues(p.Product)
}

func (p *ProductDetail) Picture(index int) string {
	return productservice.ProductPictrues(p.Product)[index]
}

func (p *ProductDetail) SupplierName(id int) string {
	if id <= 0 {
		return ""
	}
	person := personservice.GetPerson(id)
	if person != nil {
		return person.Name
	}
	return "供货商_" + string(id)
}
