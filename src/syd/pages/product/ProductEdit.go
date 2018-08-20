package product

import (
	"encoding/json"
	"fmt"
	"strings"
	"syd/model"
	"syd/service"

	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"github.com/elivoa/gxl"
)

/* ________________________________________________________________________________
   Product Create Page
*/
type ProductEdit struct {
	core.Page

	// field
	Title    string
	SubTitle string

	// property
	Id      *gxl.Int       `path-param:"1"`
	Product *model.Product `` // Product Model

	// helper used because angularjs
	Colors []*model.Object
	Sizes  []*model.Object
	Stocks []int // receive stock numbers, transfer to product later.

	Pictures []string // uploaded picture's key

	// ...
	Referer string `query:"referer"` // referer page, view or list

	// display
	StockJson string
}

func (p *ProductEdit) ProductJson() *model.Product {
	return p.Product
}

// init this page
func (p *ProductEdit) New() *ProductEdit {
	return &ProductEdit{}
}

func (p *ProductEdit) Setup() { // (string, string) {
	// page values
	p.Title = "create product post"
	fmt.Println(p.Id)
	fmt.Println(p.Product)
	if p.Id != nil {
		var err error
		if p.Product, err = service.Product.GetFullProduct(p.Id.Int); err != nil {
			panic(err)
		}

		// fill Colors and Sizes for ng.
		if p.Product != nil {
			if p.Product.Colors != nil {
				for _, c := range p.Product.Colors {
					p.Colors = append(p.Colors, model.NewObject(c))
				}
			}
			if p.Product.Sizes != nil {
				for _, c := range p.Product.Sizes {
					p.Sizes = append(p.Sizes, model.NewObject(c))
				}
			}
		}
		p.SubTitle = "编辑"
	} else {
		p.Product = model.NewProduct()
		p.SubTitle = "新建"
		p.Colors = model.NewObjectArray(2, "")
		p.Sizes = model.NewObjectArray(2, "")
	}

	// stock json
	if p.Product.Stocks != nil && len(p.Product.Stocks) > 0 {
		jsonbytes, err := json.Marshal(p.Product.Stocks)
		if err != nil {
			p.StockJson = "{}"
		}
		p.StockJson = string(jsonbytes)
		// p.StockJson = p.StockJson[1 : len(p.StockJson)-1]
	} else {
		p.StockJson = "{}"
	}
}

func (p *ProductEdit) OnPrepareForSubmitFromProductForm() {
	if p.Id == nil { // if create
		p.Product = model.NewProduct()
	} else {
		// if edit
		// for security reason, TODO security check here.
		// 读取了数据库的order是为了保证更新的时候不会丢失form中没有的数据；
		model, err := service.Product.GetFullProduct(p.Id.Int)
		if err != nil {
			panic(err.Error())
		}
		p.Product = model
		// 但是这样做就必须清除form更新的时候需要删除的值，否则form提交和原有值是叠加的，会引起错误；
		// 这里只需要清除列表等数据，这个Order中只有Details是列表。
		p.Product.ClearColors()
		p.Product.ClearSizes()
		p.Product.ClearValues()
	}
}

func (p *ProductEdit) OnSuccessFromProductForm() *exit.Exit {
	// clear values
	p.Product.ClearValues()

	// transfer stocks value to product.Stocks
	if p.Stocks != nil {

		// p.Product.Stocks = make([]*model.ProductStockItem, len(p.Product.Colors)*len(p.Product.Sizes))
		p.Product.Stocks = model.NewStocks()
		i := 0
		for _, color := range p.Product.Colors {
			for _, size := range p.Product.Sizes {
				p.Product.Stocks.Set(color, size, p.Stocks[i])
				i = i + 1
			}
		}
	}

	// transfer pictures value to pictures.
	if p.Pictures != nil {
		p.Product.Pictures = strings.Join(p.Pictures, ";")
	}

	// write to db
	if p.Id != nil {
		service.Product.UpdateProduct(p.Product)
	} else {
		service.Product.CreateProduct(p.Product)
	}

	if p.Referer == "view" {
		return exit.Redirect(fmt.Sprintf("/product/detail/%v", p.Product.Id))
	}
	// TODO: return to original page.
	return exit.Redirect("/product/list")
}
