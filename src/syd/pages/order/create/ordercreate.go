package order

import (
	"fmt"
	"got/core"
	"got/register"
	"syd/model"
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
	CustomerId int `query:"customer"`

	Customer *model.Customer
}

func (p *OrderCreateDetail) Setup() {
	p.Customer = personservice.GetCustomer(p.CustomerId)
	if p.Customer == nil {
		panic(fmt.Sprintf("customer not found: id: %v", p.CustomerId))
	}
	// Validate Privileges.
	
	
}
