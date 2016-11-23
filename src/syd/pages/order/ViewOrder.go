package order

import (
	"fmt"
	"github.com/elivoa/got/core"
	"syd/model"
	"syd/service"
	"syd/service/orderservice"
)

/* ________________________________________________________________________________
   Order Detail
*/
type ViewOrder struct {
	core.Page
	TrackNumber int64 `path-param:"1"`
	Order       *model.Order
	Customer    *model.Person
}

func (p *ViewOrder) Setup() {
	order, err := orderservice.GetOrderByTrackingNumber(p.TrackNumber)
	if err != nil {
		panic(err.Error())
	}
	p.Order = order
	p.Customer, err = service.Person.GetPersonById(p.Order.CustomerId)
	if err != nil {
		panic(err)
	}
	if p.Customer == nil {
		panic(fmt.Sprintf("customer not found: id: %v", p.Order.CustomerId))
	}
}

func (p *ViewOrder) ProductDetailJson() interface{} {
	return orderservice.OrderDetailsJson(p.Order, false)
}
