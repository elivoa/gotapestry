package product

import (
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"got/core"
	"got/register"
	"got/route"
	"gxl"
	"log"
	"net/http"
	"strconv"
	"syd/dal"
	"syd/model"
	"syd/service/productservice"
)

func Register() {
	register.Page(Register, &ProductIndex{}, &ProductEdit{}, &ProductList{},
		&ProductCreate{},
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
	Id      *gxl.Int       `path-param:"1"`
	Product *model.Product `` // Product Model
}

// init this page
func (p *ProductEdit) New() *ProductEdit {
	return &ProductEdit{}
}

func (p *ProductEdit) Setup() { // (string, string) {
	log.Println("[building] Page.SetupRender()")
	log.Println("[product] enter create/edit product")

	p.Title = "create product post"

	if p.Id != nil {
		p.Product = productservice.GetProduct(p.Id.Int)
		p.SubTitle = "编辑"
	} else {
		p.Product = model.NewProduct()
		p.SubTitle = "新建"
	}
	//return "template", "product-edit"
}

func (p *ProductEdit) OnSuccessFromProductForm() (string, string) {
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
	Products *[]model.Product
}

func (p *ProductList) Setup() {
	p.Products = dal.ListProduct()
}

// --------------------------------------------------------------------------------

var (
	//	listTypeLabel   = map[string]string{"customer": "客户", "factory": "厂商"}
	createEditLabel = map[string]string{"create": "新建", "edit": "编辑"}
)

/*
   Old handler
   -------------------------------------------------------------------------------
*/

func ProductDeletePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// TODO: important: need some security validation
	dal.DeleteProduct(id)

	// redirect to person list.
	http.Redirect(w, r, "/product/list", http.StatusFound)
}

/*
   Product Details page
   -------------------------------------------------------------------------------
*/
type ProductDetail struct {
	core.Page
	Id      int "required" // product Id
	Product *model.Product
}

func (p *ProductDetail) SetupRender() {
	id, err := strconv.Atoi(mux.Vars(p.R)["id"])
	if err == nil {
		p.Product = dal.GetProduct(id)
	}
	context.Set(p.R, route.TemplateKey, "product-detail")
}
