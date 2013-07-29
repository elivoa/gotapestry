package order

import (
	"fmt"
	"got/core"
	"got/register"
	"syd/model"
	"syd/service/orderservice"
	"syd/service/personservice"
)

func Register() {}
func init() {
	register.Component(Register, &OrderList{})
}

// ________________________________________________________________________________
//

type OrderList struct {
	core.Component

	Orders []*model.Order
	Tab    string // receive status tabs. TODO A Better way to do this?

	// temp values
	customerNames map[int]*model.Person // order-id -> customer names
}

func (p *OrderList) SetupRender() {
	// fetch customer names
	// TODO b'atch it
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

// ________________________________________________________________________________
// Events
//
func (p *OrderList) Ondelete(trackNumber int64) (string, string) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("Delete order %v \n", trackNumber)
	if _, err := orderservice.DeleteOrder(trackNumber); err != nil {
		panic(err.Error())
	}
	return "redirect", "/order/list"
}

func (p *OrderList) EditLink(order *model.Order) string {
	switch model.OrderType(order.Type) {
	case model.Wholesale:
		return fmt.Sprintf("/order/create/detail/%v", order.Id)
	case model.ShippingInstead:
		return fmt.Sprintf("/order/create/shippinginstead/%v", order.TrackNumber)
	}
	panic(fmt.Sprintf("Wrong order type for %v", order.TrackNumber))
}

func (p *OrderList) ViewLink(order *model.Order) string {
	switch model.OrderType(order.Type) {
	case model.Wholesale:
		return fmt.Sprintf("/order/view/%v", order.TrackNumber)
	case model.ShippingInstead:
		return fmt.Sprintf("/order/create/shippinginstead/%v?readonly=true", order.TrackNumber)
	}
	panic(fmt.Sprintf("Wrong order type for %v", order.TrackNumber))
}

func (p *OrderList) PrintOrderLink(order *model.Order) string {
	switch model.OrderType(order.Type) {
	case model.Wholesale:
		return fmt.Sprintf("/order/print/%v", order.TrackNumber)
	case model.ShippingInstead:
		return fmt.Sprintf("/order/shippinginsteadprint/%v", order.TrackNumber)
	}
	panic(fmt.Sprintf("Wrong order type for %v", order.TrackNumber))
}
