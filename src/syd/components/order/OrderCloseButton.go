package order

import (
	"fmt"
	"got/core"
	"syd/model"
	"syd/service/orderservice"
	"syd/service/personservice"
)

type OrderCloseButton struct {
	core.Component

	TrackNumber int64
	Source      string // return to this place

	Order    *model.Order
	Customer *model.Person
}

func (p *OrderCloseButton) Setup() {
	order, err := orderservice.GetOrderByTrackingNumber(p.TrackNumber)
	if err != nil {
		panic(err.Error())
	}
	p.Order = order
	person := personservice.GetPerson(order.CustomerId)
	if person == nil {
		panic(fmt.Sprintf("Customer not found, id: %v", order.CustomerId))
	}
	p.Customer = person
}
