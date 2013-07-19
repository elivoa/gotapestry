package order

import (
	"got/core"
	"got/register"
	"syd/model"
	"syd/service/orderservice"
	"syd/service/personservice"
	"fmt"
)

func init() {
	register.Component(Register, &OrderCloseButton{})
}

type OrderCloseButton struct {
	core.Component

	Id          string // client id
	Tid         string // component id
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
