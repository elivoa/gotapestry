package product

import (
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/db"
	"github.com/elivoa/got/route"
	"github.com/elivoa/got/route/exit"
	"strings"
	"syd/base/product"
	"syd/model"
	"syd/service"
	"time"
)

/*
   Product List page
   -------------------------------------------------------------------------------
*/
type ProductList struct {
	core.Page
	products []*model.Product
	Capital  string `path-param:"1"`
	ShowAll  bool   `query:"showall"`
	Referer  string `query:"referer"` // return here if non-empty
}

func (p *ProductList) Setup() {
}

// get main parser for prefix letter
func (p *ProductList) getMainParser(letter string) *db.QueryParser { // []*model.Product
	var parser = service.Order.EntityManager().NewQueryParser()
	p.Capital = strings.ToLower(p.Capital)
	if p.Capital == "" || p.Capital == "all" {
		parser.Where().Limit(1000) // disable default limit
	} else {
		parser.Where("capital", p.Capital)
	}
	return parser
}

func (p *ProductList) Products(letter string) []*model.Product {
	if nil != p.products {
		return p.products
	}

	parser := p.getMainParser(letter)
	products, err := service.Product.List(parser, service.WITH_PRODUCT_DETAIL|service.WITH_PRODUCT_INVENTORY)
	// products, err := service.Product.List(parser, service.WITH_NONE)
	if nil != err {
		panic(err.Error())
	}
	p.products = products
	return products
}

// data ajax entrance.
func (p *ProductList) Ongetproducts(letter string) *exit.Exit {
	products := p.Products(letter)
	return exit.MarshalJson(copyToBasicProduct(products))
}

// data ajax: products.
func (p *ProductList) Ongetproductstocks(letter string) *exit.Exit {
	products := p.Products(letter)
	return exit.MarshalJson(copyToProductStocks(products))
}

// data ajax: products.
func (p *ProductList) Ongetproductdetails(letter string) *exit.Exit {
	products := p.Products(letter)
	return exit.MarshalJson(copyToProductDetails(products))
}

// NOTE: event name is case sensitive. Kill this when add cache.
func (p *ProductList) Ondelete(productId int) *exit.Exit {
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

// --------------------------------------------------------------------------------
// for output json

type ProductListJsonObject struct {
	Id           int            // id
	Name         string         // product name
	ProductId    string         // 传说中的货号
	Status       product.Status //
	Brand        string         `json:",omitempty"`
	Price        float64        `json:",omitempty"`
	Supplier     int            `json:"-"`
	FactoryPrice float64        `json:"-"`
	Stock        int            `json:"-"`          // 库存量 || not used again?
	ShelfNo      string         `json:"-"`          // 货架号
	Capital      string         `json:",omitempty"` // captical letter to quick access.
	Note         string         `json:",omitempty"`
	CreateTime   time.Time      `json:"-"`
	UpdateTime   time.Time      `json:"-"`
}

type ProductStocksJsonObject struct {
	Id     int
	Stock  int          `json:",omitempty"`
	Stocks model.Stocks `json:",omitempty"` // map[string]int
}

type ProductDetailsJsonObject struct {
	Id         int
	Colors     []string            `json:",omitempty"` // stores in product_properties table.
	Sizes      []string            `json:",omitempty"`
	Properties map[string][]string `json:",omitempty"` // other properties // TODO
}

func copyToBasicProduct(products []*model.Product) []*ProductListJsonObject {
	basicps := []*ProductListJsonObject{}
	if nil != products {
		for _, p := range products {
			basicp := &ProductListJsonObject{
				Id:           p.Id,
				Name:         p.Name,
				ProductId:    p.ProductId,
				Status:       p.Status,
				Brand:        p.Brand,
				Price:        p.Price,
				Supplier:     p.Supplier,
				FactoryPrice: p.FactoryPrice,
				Stock:        p.Stock,
				ShelfNo:      p.ShelfNo,
				Capital:      p.Capital,
				Note:         p.Note,
				CreateTime:   p.CreateTime,
				UpdateTime:   p.UpdateTime,
			}
			basicps = append(basicps, basicp)
		}
	}
	return basicps
}

func copyToProductStocks(products []*model.Product) []*ProductStocksJsonObject {
	basicps := []*ProductStocksJsonObject{}
	if nil != products {
		for _, p := range products {
			basicp := &ProductStocksJsonObject{
				Id:     p.Id,
				Stock:  p.Stock,
				Stocks: p.Stocks,
			}
			basicps = append(basicps, basicp)
		}
	}
	return basicps
}

func copyToProductDetails(products []*model.Product) []*ProductDetailsJsonObject {
	basicps := []*ProductDetailsJsonObject{}
	if nil != products {
		for _, p := range products {
			basicp := &ProductDetailsJsonObject{
				Id:         p.Id,
				Colors:     p.Colors,
				Sizes:      p.Sizes,
				Properties: p.Properties,
			}
			basicps = append(basicps, basicp)
		}
	}
	return basicps
}
