package product

import (
	"got/core"
	"strings"
	"syd/model"
	"syd/service/productservice"
)

/* ________________________________________________________________________________
   Product Home Page
*/
type ProductIndex struct{ core.Page }

func (p *ProductIndex) Setup() (string, string) { return "redirect", "/product/list" }

// redirect
type ProductCreate struct{ core.Page }

func (p *ProductCreate) Setup() (string, string) { return "redirect", "/product/edit" }

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
	p.Capital = strings.ToLower(p.Capital)
	if p.Capital == "" || p.Capital == "all" {
		p.Products, err = productservice.ListProducts()
	} else {
		p.Products, err = productservice.ListProductsByCapital(p.Capital)
	}
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

// --------------------------------------------------------------------------------

var (
	//	listTypeLabel   = map[string]string{"customer": "客户", "factory": "厂商"}
	createEditLabel = map[string]string{"create": "新建", "edit": "编辑"}
)
