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

	// Validate Privileges.
	if p.Id != nil {
		// edit mode
		order, err := orderservice.GetOrder(p.Id.Int)
		if err != nil {
			panic(err.Error())
		}
		p.Order = order
		p.CustomerId = p.Order.CustomerId

	} else {
		// create mode
		p.Order = model.NewOrder()
		p.Order.CustomerId = p.CustomerId
	}

	// init person
	p.Customer = personservice.GetCustomer(p.CustomerId)
	if p.Customer == nil {
		panic(fmt.Sprintf("customer not found: id: %v", p.CustomerId))
	}
}

// before submit
func (p *OrderCreateDetail) OnSubmit() {
	if p.Id == nil {
		// if create
		p.Order = model.NewOrder()
		p.Order.CustomerId = p.CustomerId
	} else {
		// if edit
		// for security reason, TODO security check here.
		o, err := orderservice.GetOrder(p.Id.Int)
		if err != nil {
			panic(err.Error())
		}
		p.Order = o
	}
}

// after submit
func (p *OrderCreateDetail) OnSuccess() (string, string) {
	// order
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
