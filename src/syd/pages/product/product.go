package product

import (
	"encoding/json"
	"fmt"
	"got/core"
	"got/debug"
	"got/register"
	"gxl"
	"strings"
	"syd/dal"
	"syd/model"
	"syd/service/personservice"
	"syd/service/productservice"
)

func Register() {
	register.Page(Register,
		&ProductIndex{}, &ProductEdit{}, &ProductList{},
		&ProductCreate{}, &ProductDetail{},
	)
}

/* ________________________________________________________________________________
   Product Home Page
*/
type ProductIndex struct{ core.Page }

func (p *ProductIndex) Setup() (string, string) { return "redirect", "/product/list" }

// redirect
type ProductCreate struct{ core.Page }

func (p *ProductCreate) Setup() (string, string) { return "redirect", "/product/edit" }

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
	if p.Id != nil {
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

func (p *ProductEdit) OnSuccessFromProductForm() (string, string) {
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
	return "redirect", "/product/list"
}

/*
   Product List page
   -------------------------------------------------------------------------------
*/
type ProductList struct {
	core.Page
	Products []*model.Product
}

func (p *ProductList) Setup() {
	var err error
	p.Products, err = productservice.ListProducts()
	if nil != err {
		panic(err.Error())
		// Goto error page
	}
	// p.Products = dal.ListProduct()
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

// NOTE: event name is case sensitive. Kill this when add cache.
func (p *ProductList) Ondelete(productId int) (string, string) {
	debug.Log("Delete Product %d", productId)
	dal.DeleteProduct(productId)
	// TODO make this default redirect.
	return "redirect", "/product/list"
}

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

// --------------------------------------------------------------------------------

var (
	//	listTypeLabel   = map[string]string{"customer": "客户", "factory": "厂商"}
	createEditLabel = map[string]string{"create": "新建", "edit": "编辑"}
)
