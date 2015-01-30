package product

// Deprecated, TODO chagne this into angularjs module.
import (
	"bytes"
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
	if nil == product {
		return template.HTML("ERROR: PRODUCT IS NIL!")
	}

	var spec bytes.Buffer
	if product.Colors != nil {
		for idx, color := range product.Colors {
			if idx > 0 {
				spec.WriteString(", ")
			}
			spec.WriteString(color)
		}
	}
	if product.Sizes != nil {
		spec.WriteString("<span class=\"vline\">|</span>")
		for idx, size := range product.Sizes {
			if idx > 0 {
				spec.WriteString(", ")
			}
			spec.WriteString(size)
		}
	}
	return template.HTML(spec.String())
}

func (p *ProductList) StockDescription(product *model.Product) template.HTML {
	if nil != product {
		if nil != product.Stocks && len(product.Stocks) > 0 {
			var buffer bytes.Buffer

			product.Stocks.Loop(func(color, size string, stock int) {
				// inventorydao.SetProductStock(product.Id, color, size, stock)
				buffer.WriteString(color)
				buffer.WriteString("/")
				buffer.WriteString(size)
				buffer.WriteString(":")
				buffer.WriteString(strconv.Itoa(stock))
				buffer.WriteString(";&#10; ")
			})
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
