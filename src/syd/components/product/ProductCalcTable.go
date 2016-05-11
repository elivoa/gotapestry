package product

// Deprecated, TODO chagne this into angularjs module.
import (
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"syd/model"
	"syd/service"
)

type ProductCalcTable struct {
	core.Component

	Data     *model.ProductSalesTable
	products map[int64]*model.Product
}

// display: total stocks
func (p *ProductCalcTable) Setup() *exit.Exit {

	if nil == p.Data {
		panic("Parameter Data Must be set;")
	}

	// get product
	productmap, err := service.Product.BatchFetchProductByIdMap(p.Data.ProductMap())
	if err != nil {
		return exit.Error(err)
	}
	p.products = productmap

	return nil
}

func (p *ProductCalcTable) GetData(date string, productId int64) string {
	value := p.Data.Get(date, productId)
	if value == 0 {
		return ""
	} else {
		return fmt.Sprint(value)
	}
}

var emptyproduct = model.NewProduct()

func (p *ProductCalcTable) GetProduct(productId int64) *model.Product {
	if nil != p.products {
		if product, ok := p.products[productId]; ok && nil != product {
			return product
		}
	}
	return emptyproduct
}
