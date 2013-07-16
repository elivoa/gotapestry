package order

import (
	"fmt"
	"github.com/gorilla/mux"
	"got/core"
	"got/debug"
	"got/register"
	"got/route"
	"got/templates"
	"gxl"
	"strings"
	"syd/dal"
	"syd/dal/productdao"
	"syd/model"
	"syd/service/orderservice"
)

/* ________________________________________________________________________________
   Register all pages under /order
*/
func init() {
	register.Page(Register, &OrderList{}, &OrderIndex{}, &OrderEdit{})
}
func Register() {}

// ________________________________________________________________________________
// Start Order pages
//
// ____ Order Index _______________________________________________________________
type OrderIndex struct {
	core.Page           `PageRedirect:"/order/list"`
	__got_page_redirect int `PageRedirect:"/order/list"` //?
}

func (p *OrderIndex) SetupRender() (string, string) {
	return "redirect", "/order/list"
}

/* ________________________________________________________________________________
   Order List
*/
type OrderList struct {
	core.Page

	Orders []*model.Order
	Tab    string `path-param:"1"`

	// customerNames map[int]*model.Person // order-id -> customer names
}

func (p *OrderList) Activate() {
	if p.Tab == "" {
		p.Tab = "all"
	}
}

func (p *OrderList) SetupRender() {
	orders, err := orderservice.ListOrder(p.Tab)
	if err != nil {
		debug.Error(err)
		panic(err.Error())
	}
	p.Orders = orders
	// p.Orders = dal.ListOrder(p.Tab)
}

func (p *OrderList) TabStyle(tab string) string {
	if strings.ToLower(p.Tab) == strings.ToLower(tab) {
		return "cur"
	}
	return ""
}

/* ________________________________________________________________________________
   Order Detail
*/
type OrderDetail struct {
	core.Page
}

/* ________________________________________________________________________________
   Order Create/Edit
*/
func New() *OrderPage {
	// templates.Add("order-edit")
	// templates.Add("order-list")
	templates.Add("order-detail")
	return &OrderPage{}
}

type OrderPage struct{}

func (p *OrderPage) Mapping(r *mux.Router) {
	// r.HandleFunc("/order/", route.RedirectHandler("/order/list"))
	// r.HandleFunc("/order/list", route.PageHandler(&OrderList{}))
	r.HandleFunc("/order/{id:[0-9]+}", route.PageHandler(&OrderDetail{}))

	// editHandler := route.PageHandler(&OrderEdit{})
	// r.HandleFunc("/order/create", editHandler)
	// r.HandleFunc("/order/edit/{id:[0-9]+}", editHandler)
}

/* ________________________________________________________________________________
   Order Edit
*/
type OrderEdit struct {
	core.Page
	Title       string // `path-param:"2"`
	SubTitle    string
	SubmitLabel string

	Id    *gxl.Int     `path-param:"1"` // is this param useful?
	Order *model.Order ``
}

// func (p *OrderEdit) Activate(id int64) {
// 	p.Id = big.NewInt(id)
// }

func (p *OrderEdit) SetupRender() (interface{}, string) {
	var err error
	if p.Id != nil {
		p.Order, err = dal.GetOrder(p.Id.Int)
		if err != nil {
			return err, ""
		}
		p.Title = "编辑订单"
		p.SubTitle = "编辑"
		p.SubmitLabel = "保存"
	} else {
		p.Order = model.NewOrder()
		p.Title = "新建订单"
		p.SubTitle = "新建"
		p.SubmitLabel = "保存"
	}
	return nil, ""
}

func (p *OrderEdit) ProductDisplayName(id int) string {
	product, _ := productdao.Get(id)
	if product != nil {
		return product.Name
	}
	return "【error】"
}

// func (p *OrderEdit) AfterRender() (interface{}, string) {
// 	return "template", "order-edit"
// }

func (p *OrderEdit) OnSuccessFromOrderForm() (string, string) {

	// debug log
	fmt.Printf("+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+\n")
	fmt.Printf("+ OrderDetails %v\n", len(p.Order.Details))
	for idx, o := range p.Order.Details {
		fmt.Printf("+ %v: %v\n", idx, o)
	}
	fmt.Println()

	// modify logic
	switch p.Order.DeliveryMethod {
	case "Express":
		p.Order.Status = "New"
	case "TakeAway":
		p.Order.Status = "Delivering"
	}

	// update
	if p.Id != nil {
		dal.UpdateOrder(p.Order)
	} else {
		dal.CreateOrder(p.Order)
	}
	return "redirect", "/order/list"
}
