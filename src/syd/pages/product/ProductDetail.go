package product

import (
	"github.com/elivoa/got/core"
	"github.com/elivoa/gxl"
	"syd/model"
	"syd/service"
)

// ________________________________________________________________________________
// Product Details

type ProductDetail struct {
	core.Page
	Id      *gxl.Int `path-param:"1"`
	Product *model.Product
}

func (p *ProductDetail) Setup() {
	var err error
	p.Product, err = service.Product.GetProduct(p.Id.Int)
	if err != nil {
		panic(err)
	}
}

func (p *ProductDetail) Pictures() []string {
	return service.Product.ProductPictrues(p.Product)
}

func (p *ProductDetail) Picture(index int) string {
	return service.Product.ProductPictrues(p.Product)[index]
}

func (p *ProductDetail) SupplierName(id int) string {
	if id <= 0 {
		return ""
	}
	if person, err := service.Person.GetPersonById(id); err != nil {
		panic(err)
	} else if person != nil {
		return person.Name
	}
	return "供货商_" + string(id)
}
