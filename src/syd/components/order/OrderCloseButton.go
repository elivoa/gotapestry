package order

import (
	"github.com/elivoa/got/core"
	"syd/model"
)

type OrderCloseButton struct {
	core.Component
	Order *model.Order

	// TrackNumber int64
	// Customer    *model.Person
	// Referer     string
}

func (p *OrderCloseButton) Setup() {
	// if p.Order == nil {
	// 	order, err := orderservice.GetOrderByTrackingNumber(p.TrackNumber)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	p.Order = order
	// } else {
	// 	// set customer if nil
	// 	if p.Customer == nil && p.Order.Customer != nil {
	// 		p.Customer = p.Order.Customer
	// 	}
	// 	p.TrackNumber = p.Order.TrackNumber
	// }

	// if p.Customer == nil {

	// 	person, err := service.Person.GetPersonById(p.Order.CustomerId)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if person == nil {
	// 		panic(fmt.Sprintf("Customer not found, id: %v", p.Order.CustomerId))
	// 	}
	// 	p.Customer = person
	// }
}
