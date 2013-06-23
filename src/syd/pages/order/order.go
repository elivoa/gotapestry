package order

import (
	"fmt"
	"github.com/gorilla/mux"
	"got/core"
	"got/register"
	"got/route"
	"got/templates"
	"gxl"
	"syd/dal"
	"syd/model"
)

/* ________________________________________________________________________________
   Got Style page
*/
func Register() {}
func init() {
	register.Page(Register, &OrderList{}, &OrderIndex{}, &OrderEdit{})
}

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
	Orders *[]model.Order

	customerNames map[int]*model.Person // order-id -> customer names
}

func (p *OrderList) SetupRender() {
	p.Orders = dal.ListOrder("all")

	// fetch customer names
	// TODO batch it
	length := len(*p.Orders)
	if length > 0 {
		p.customerNames = make(map[int]*model.Person, length)
		for _, o := range *p.Orders {
			if _, ok := p.customerNames[o.CustomerId]; ok {
				continue
			}

			customer := dal.GetPerson(o.CustomerId)
			if customer != nil {
				p.customerNames[customer.Id] = customer
			}
		}
	}

	// return "template", "order-list"
}

func (p *OrderList) ShowCustomerName(customerId int) string {
	customer, ok := p.customerNames[customerId]
	if ok {
		return customer.Name
	} else {
		return fmt.Sprintf("_[ p%v ]_", customerId)
	}
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
	Title       string `path-param:"2"`
	SubTitle    string
	SubmitLabel string

	Id    *gxl.Int     `param:"id" path-param:"1"` // is this param useful?
	Order *model.Order ``
}

// func (p *OrderEdit) Activate(id int64) {
// 	p.Id = big.NewInt(id)
// }

func (p *OrderEdit) SetupRender() (interface{}, string) {
	// fmt.Printf("============================-----=====-----===== %v\n", p.Id)
	// fmt.Printf("============================-----=====-----===== %v\n", p.Title)

	var err error
	if p.Injected("Id") && p.Id != nil {
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
	product := dal.GetProduct(id)
	if product != nil {
		return product.Name
	}
	return "【error】"
}

// func (p *OrderEdit) AfterRender() (interface{}, string) {
// 	return "template", "order-edit"
// }

func (p *OrderEdit) OnSuccessFromOrderForm() (string, string) {

	fmt.Printf("+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+\n")
	fmt.Printf("+ OrderDetails %v\n", len(p.Order.Details))
	for idx, o := range p.Order.Details {
		fmt.Printf("+ %v: %v\n", idx, o)
	}

	if p.Injected("Id") && p.Id != nil {
		dal.UpdateOrder(p.Order)
	} else {
		dal.CreateOrder(p.Order)
	}

	// _, err := strconv.Atoi(mux.Vars(p.R)["id"])
	// if err == nil {
	// } else {
	// }
	return "redirect", "/order/list"
}
