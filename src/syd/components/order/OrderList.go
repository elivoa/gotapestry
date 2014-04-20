package order

import (
	"fmt"
	"got/core"
	"syd/model"
	"syd/service"
	"syd/service/orderservice"
	"syd/service/personservice"
)

type OrderList struct {
	core.Component

	Orders     []*model.Order
	Tab        string  // receive status tabs. TODO A Better way to do this?
	TotalPrice float64 // all order's price

	// temp values
	customerNames map[int]*model.Person // order-id -> customer names
}

func (p *OrderList) SetupRender() {
	// verify user role.
	service.User.RequireRole(p.W, p.R, "admin") // TODO remove w, r. use service injection.

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
			p.TotalPrice += o.TotalPrice
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
func (p *OrderList) Ondelete(trackNumber int64, tab string) (string, string) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("Delete order %v \n", trackNumber)
	if _, err := orderservice.DeleteOrder(trackNumber); err != nil {
		panic(err)
	}
	return "redirect", fmt.Sprintf("/order/list/%v", tab)
}

// OnPrint set order's status to `todeliver` then go to print page.
func (p *OrderList) OnPrint(trackNumber int64) (string, string) {
	if err := orderservice.ChangeOrderStatus(trackNumber, "todeliver"); err != nil {
		panic(err.Error())
	}
	return "redirect", fmt.Sprintf("/order/print/%v", trackNumber)
}

// shipping instead order's status changed to delivering
func (p *OrderList) OnShippingInsteadOrderPrint(trackNumber int64) (string, string) {
	if err := orderservice.ChangeOrderStatus(trackNumber, "delivering"); err != nil {
		panic(err)
	}
	return "redirect", fmt.Sprintf("/order/shippinginsteadprint/%v", trackNumber)
}

func (p *OrderList) EditLink(order *model.Order) string {
	switch model.OrderType(order.Type) {
	case model.Wholesale:
		return fmt.Sprintf("/order/create/detail/%v", order.Id)
	case model.ShippingInstead:
		return fmt.Sprintf("/order/create/shippinginstead/%v", order.TrackNumber)
	case model.SubOrder:
		return "#Error"
	default:
		return fmt.Sprintf("%v", order.Type)
	}
	panic(fmt.Sprintf("Wrong order type for %v", order.TrackNumber))
}

func (p *OrderList) ViewLink(order *model.Order) string {
	switch model.OrderType(order.Type) {
	case model.Wholesale:
		return fmt.Sprintf("/order/view/%v", order.TrackNumber)
	case model.ShippingInstead:
		return fmt.Sprintf("/order/create/shippinginstead/%v?readonly=true", order.TrackNumber)
	case model.SubOrder:
		return "#Error"
	default:
		return fmt.Sprintf("%v", order.Type)
	}
	panic(fmt.Sprintf("Wrong order type for %v", order.TrackNumber))
}

// to some action and then redirect to print page.
func (p *OrderList) PrintOrderLink(order *model.Order) string {
	switch model.OrderType(order.Type) {
	case model.Wholesale:
		// TODO auto generate this via builtin eventlink component.
		return fmt.Sprintf("/order/list.orderlist.Print/%v", order.TrackNumber)
	case model.ShippingInstead:
		return fmt.Sprintf("/order/list.orderlist.ShippingInsteadOrderPrint/%v", order.TrackNumber)
	case model.SubOrder:
		return "#Error"
	default:
		return fmt.Sprintf("%v", order.Type)
	}
	panic(fmt.Sprintf("Wrong order type for %v", order.TrackNumber))
}

// Print only. do nothing.
func (p *OrderList) FixPrintOrderLink(order *model.Order) string {
	switch model.OrderType(order.Type) {
	case model.Wholesale:
		return fmt.Sprintf("/order/print/%v", order.TrackNumber)
	case model.ShippingInstead:
		return fmt.Sprintf("/order/shippinginsteadprint/%v", order.TrackNumber)
	case model.SubOrder:
		return "#Error"
	default:
		return fmt.Sprintf("%v", order.Type)
	}
	panic(fmt.Sprintf("Wrong order type for %v", order.TrackNumber))
}
