package product

import (
	"encoding/json"
	"fmt"
	"github.com/elivoa/got/route/exit"
	"github.com/elivoa/gxl"
	"github.com/elivoa/got/core"
	"strings"
	"syd/model"
	"syd/service/productservice"
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
	Id       *gxl.Int       `path-param:"1"`
	Product  *model.Product `` // Product Model
	Stocks   []int          // receive stock numbers, transfer to product later.
	Pictures []string       // uploaded picture's key

	Referer string `query:"referer"` // referer page, view or list

	// display
	StockJson string
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
		fmt.Printf("\t >>> get product by id\n")

		p.Product = productservice.GetProduct(p.Id.Int)
		p.SubTitle = "编辑"
	} else {
		p.Product = model.NewProduct()
		p.SubTitle = "新建"
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

func (p *ProductEdit) OnSuccessFromProductForm() *exit.Exit {
	// clear values
	p.Product.ClearValues()

	// transfer stocks value to product.Stocks
	if p.Stocks != nil {
		p.Product.Stocks = map[string]int{}
		i := 0
		for _, color := range p.Product.Colors {
			for _, size := range p.Product.Sizes {
				key := fmt.Sprintf("%v__%v", color, size)
				p.Product.Stocks[key] = p.Stocks[i]
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
		productservice.UpdateProduct(p.Product)
	} else {
		productservice.CreateProduct(p.Product)
	}

	if p.Referer == "view" {
		return exit.Redirect(fmt.Sprintf("/product/detail/%v", p.Product.Id))
	}
	return exit.Redirect("/product/list")
	// return "redirect", "/product/list"
}
