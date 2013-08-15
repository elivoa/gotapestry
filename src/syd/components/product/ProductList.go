package product

import (
	"bytes"
	"got/core"
	"got/debug"
	"got/route"
	"html/template"
	"syd/dal"
	"syd/model"
	"syd/service/productservice"
)

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

func (p *ProductList) ShowSpecification(productId int) template.HTML {
	product := productservice.GetProduct(productId)
	if nil == product{
		return template.HTML("ERROR: PRODUCT IS NIL!")
	}
	var spec bytes.Buffer
	if product.Colors == nil || len(product.Colors) > 0 {
		i := 0
		for _, color := range product.Colors {
			if i = i + 1; i > 1 {
				spec.WriteString("<span class=\"vline\">|</span>")
			}
			spec.WriteString(color)
			spec.WriteString(": ")
			o := 0
			for _, size := range product.Sizes {
				// if stock is nil, returns nothing.
				// to be continued....
				if o = o + 1; o > 1 {
					spec.WriteString(" / ")
				}
				spec.WriteString(size)
			}
		}
	}
	return template.HTML(spec.String())
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
