package order

import (
	"fmt"
	"github.com/elivoa/got/core"
	"syd/model"
	"syd/service"
	"syd/service/orderservice"
)

type OrderCloseButton struct {
	core.Component

	TrackNumber int64
	Order       *model.Order
	Customer    *model.Person
	Referer     string
}

func (p *OrderCloseButton) Setup() {
	order, err := orderservice.GetOrderByTrackingNumber(p.TrackNumber)
	if err != nil {
		panic(err)
	}
	p.Order = order
	person, err := service.Person.GetPersonById(order.CustomerId)
	if err != nil {
		panic(err)
	}
	if person == nil {
		panic(fmt.Sprintf("Customer not found, id: %v", order.CustomerId))
	}
	p.Customer = person
}
