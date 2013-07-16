package order

import (
	"fmt"
	"got/core"
	"got/register"
	"syd/model"
	"syd/service/personservice"
)

func Register() {}
func init() {
	register.Component(Register, &OrderList{})
}

type OrderList struct {
	core.Component

	Orders []*model.Order

	// temp values
	customerNames map[int]*model.Person // order-id -> customer names
}

func (p *OrderList) SetupRender() {
	// fetch customer names
	// TODO batch it
	if p.Orders == nil {
		return
	}

	// Prepare customerNames to display.
	length := len(p.Orders)
	if length > 0 {
		p.customerNames = make(map[int]*model.Person, length)
		for _, o := range p.Orders {
			if _, ok := p.customerNames[o.CustomerId]; ok {
				continue
			}

			customer := personservice.GetPerson(o.CustomerId)
			if customer != nil {
				p.customerNames[customer.Id] = customer
			}
		}
	}
}

func (p *OrderList) ShowCustomerName(customerId int) string {
	customer, ok := p.customerNames[customerId]
	if ok {
		return customer.Name
	} else {
		return fmt.Sprintf("_[ p%v ]_", customerId)
	}
}
