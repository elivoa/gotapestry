package product

import (
	"bytes"
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/debug"
	"github.com/elivoa/got/route"
	"github.com/elivoa/got/route/exit"
	"html/template"
	"strconv"
	"syd/dal/inventorydao"
	"syd/model"
	"syd/service"
)

type ProductList struct {
	core.Component
	Products []*model.Product
	Source   string `query:"source"` // return here
}

// NOTE: event name is case sensitive. Kill this when add cache.
func (p *ProductList) Ondelete(productId int) *exit.Exit {
	debug.Log("Delete Product %d", productId)
	service.Product.DeleteProduct(productId)
	// TODO make this default redirect.
	return route.RedirectDispatch(p.Source, "/product/list")
}

// returns a string contains available size and color.
func (p *ProductList) ShowSpecification(product *model.Product) template.HTML {
	// product, err := service.Product.GetProduct(productId)
	// if err != nil {
	// 	return template.HTML(err.Error())
	// }
	if nil == product {
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
				key := fmt.Sprintf("%v__%v", color, size)
				stock := product.Stocks[key]
				if stock != -1 { // stock = -1 means not available.
					if o = o + 1; o > 1 {
						spec.WriteString(" / ")
					}
					spec.WriteString(size)
				}
			}
		}
	}
	return template.HTML(spec.String())
}

func (p *ProductList) StockDescription(product *model.Product) template.HTML {
	if nil != product {
		if nil != product.Stocks && len(product.Stocks) > 0 {
			var buffer bytes.Buffer
			for key, s := range product.Stocks {
				buffer.WriteString(key)
				buffer.WriteString(":")
				buffer.WriteString(strconv.Itoa(s))
				buffer.WriteString(";&#10; ")
			}
		}
	}
	return "没有库存或尚未清点库存！"
}

// display: total stocks
func (p *ProductList) NStock(productId int) (sum int) {
	stockmap := inventorydao.ListProductStocks(productId)
	if stockmap != nil {
		for _, stock := range stockmap {
			if stock > 0 {
				sum += stock
			}
		}
	}
	return
}
