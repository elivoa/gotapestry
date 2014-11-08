package product

import (
	"github.com/elivoa/got/core"
	"strings"
	"syd/model"
	"syd/service"
)

/*
   Product List page
   -------------------------------------------------------------------------------
*/
type ProductList struct {
	core.Page
	Products []*model.Product
	Capital  string `path-param:"1"`
}

func (p *ProductList) Setup() {
	var err error
	var parser = service.Order.EntityManager().NewQueryParser()
	p.Capital = strings.ToLower(p.Capital)
	if p.Capital == "" || p.Capital == "all" {
		parser.Where().Limit(1000) // disable default limit
	} else {
		parser.Where("capital", p.Capital)
	}
	p.Products, err = service.Product.List(parser, service.WITH_PRODUCT_DETAIL|service.WITH_PRODUCT_INVENTORY)
	if nil != err {
		panic(err.Error())
	}
}

func (p *ProductList) TabClass(letter string) string {
	if "all" == letter && p.Capital == "" {
		return "cur"
	}
	if strings.ToLower(p.Capital) == strings.ToLower(letter) {
		return "cur"
	}
	return ""
}
