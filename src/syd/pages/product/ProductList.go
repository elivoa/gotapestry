package product

import (
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/debug"
	"github.com/elivoa/got/route"
	"github.com/elivoa/got/route/exit"
	"strings"
	"syd/base/product"
	"syd/model"
	"syd/service"
)

/*
   Product List page
   -------------------------------------------------------------------------------
*/
type ProductList struct {
	core.Page
	// Products []*model.Product
	Capital string `path-param:"1"`
	ShowAll bool   `query:"showall"`
	Referer string `query:"referer"` // return here if non-empty
}

func (p *ProductList) Setup() {
}

// json method
func (p *ProductList) Products(letter string) []*model.Product {
	fmt.Println("\n ----- , ", letter)

	var parser = service.Order.EntityManager().NewQueryParser()
	p.Capital = strings.ToLower(p.Capital)
	if p.Capital == "" || p.Capital == "all" {
		parser.Where().Limit(1000) // disable default limit
	} else {
		parser.Where("capital", p.Capital)
	}
	products, err := service.Product.List(parser, service.WITH_PRODUCT_DETAIL|service.WITH_PRODUCT_INVENTORY)
	if nil != err {
		panic(err.Error())
	}
	return products
}

// NOTE: event name is case sensitive. Kill this when add cache.
func (p *ProductList) Ondelete(productId int) *exit.Exit {
	debug.Log("Delete Product xxx %d", productId)
	service.Product.DeleteProduct(productId)
	return route.RedirectDispatch(p.Referer, "product/list/"+p.Capital)
}

func (p *ProductList) Onshow(productId int) *exit.Exit {
	service.Product.ChangeStatus(productId, product.StatusNormal)
	return route.RedirectDispatch(p.Referer, "product/list/"+p.Capital)
}

func (p *ProductList) Onhide(productId int) *exit.Exit {
	service.Product.ChangeStatus(productId, product.StatusHide)
	return route.RedirectDispatch(p.Referer, "product/list/"+p.Capital)
}
