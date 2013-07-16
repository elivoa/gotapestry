package order

import (
	"fmt"
	"got/core"
	"got/register"
	"gxl"
	"syd/model"
	"syd/service/orderservice"
	"syd/service/personservice"
)

/* ________________________________________________________________________________
   Register all pages under /order
*/
func init() {
	register.Page(Register,
		&OrderCreateIndex{},
		&OrderCreateDetail{},
	)
}
func Register() {}

// ____ Order Create Index ___________________________________________________________
// BUG: OrderCreate here can't fallback to /order/create
//
type OrderCreateIndex struct {
	core.Page
	CustomerId int `param:"."`
}

func (p *OrderCreateIndex) Setup() {
	// enter person create order.
}

func (p *OrderCreateIndex) OnSuccessFromCustomerForm() (string, string) {
	if p.CustomerId > 0 {
		url := fmt.Sprintf("/order/create/detail?customer=%v", p.CustomerId)
		return "redirect", url
	}
	return "redirect", "thispage"
}

// ____ Order Create Details _______________________________________________________________
//
type OrderCreateDetail struct {
	core.Page
	CustomerId int      `query:"customer"`
	Id         *gxl.Int `path-param:"1"` // order id

	Customer *model.Customer
	Order    *model.Order // submit the big table.
	DaoFu    string       // on

	SourceUrl string `query:"source"` // redirect url
}

func (p *OrderCreateDetail) Setup() {
	if p.Id == nil && p.CustomerId == 0 {
		panic("Can't find order to edit!")
	}

	// V121alidate Privileges.
	if p.Id != nil {
		order, err := orderservice.GetOrder(p.Id.Int)
		if err != nil {
			panic(err.Error())
		}
		p.Order = order
		p.CustomerId = p.Order.CustomerId
	} else {
		p.Order = model.NewOrder()
	}

	// init person
	p.Customer = personservice.GetCustomer(p.CustomerId)
	if p.Customer == nil {
		panic(fmt.Sprintf("customer not found: id: %v", p.CustomerId))
	}

}

// before submit
func (p *OrderCreateDetail) OnSubmit() {
	p.Order = model.NewOrder()
}

// after submit
func (p *OrderCreateDetail) OnSuccess() (string, string) {
	fmt.Println("********************************************************************************")
	fmt.Println("on order form sugmit")
	fmt.Println(p.Order.DeliveryMethod)
	fmt.Println(p.Order.ExpressFee)
	fmt.Println(p.DaoFu)

	// order
	p.Order.CustomerId = p.CustomerId
	for _, detail := range p.Order.Details {
		detail.OrderTrackNumber = p.Order.TrackNumber
	}
	if p.DaoFu == "on" {
		p.Order.ExpressFee = -1
	}

	// update
	if p.Id != nil {
		orderservice.UpdateOrder(p.Order)
	} else {
		orderservice.CreateOrder(p.Order)
	}

	// return source?
	if p.SourceUrl == "" {
		return "redirect", fmt.Sprintf("/order/print/%v", p.Order.TrackNumber)
	} else {
		return "redirect", p.SourceUrl
	}
}

func (p *OrderCreateDetail) IsEdit() bool {
	return p.Id != nil
}

func (p *OrderCreateDetail) ProductDetailJson() interface{} {
	return orderservice.ProductDetailJson(p.Order)
}
