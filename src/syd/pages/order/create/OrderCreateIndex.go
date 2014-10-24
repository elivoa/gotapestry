package order

import (
	"fmt"
	"github.com/elivoa/got/core"
	"syd/model"
	"syd/service/orderservice"
)

/* ________________________________________________________________________________
   Register all pages under /order
*/
// --------  Order Create Index  -----------------------------------------------------------------------

// BUG: OrderCreate here can't fallback to /order/create
type OrderCreateIndex struct {
	core.Page
	CustomerId int    `param:"."`
	OrderType  string // wholesale | shippinginstead
}

func (p *OrderCreateIndex) Setup() {
	// enter person create order.
}

func (p *OrderCreateIndex) OnSuccessFromCustomerForm() (string, string) {
	var url string
	if p.CustomerId > 0 {
		if p.OrderType == "shippinginstead" {
			// create an order and goto ShippingInstead edit page.
			order := model.NewOrder()
			order.Type = uint(model.ShippingInstead)
			order.CustomerId = p.CustomerId
			orderservice.CreateOrder(order)
			url = fmt.Sprintf("/order/create/shippinginstead/%v", order.TrackNumber)
		} else {
			url = fmt.Sprintf("/order/create/detail?customer=%v", p.CustomerId)
		}
		return "redirect", url
	}
	return "redirect", "thispage"
}
