package orderreturns

import (
	"fmt"
	"syd/model"
	"syd/service"
	"syd/service/orderservice"

	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route"
	"github.com/elivoa/got/route/exit"
)

type OrderList struct {
	core.Component

	Orders     []*model.Order
	Tab        string // receive status tabs. TODO A Better way to do this?
	TotalItems int
	TotalPrice float64 // all order's price
	Referer    string  // return here.

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
	p.TotalItems = 0
	p.TotalPrice = 0
	if length > 0 {
		for _, o := range p.Orders {
			p.TotalPrice += o.TotalPrice
			p.TotalItems += o.TotalCount
		}
	}
}

// ________________________________________________________________________________
// Events
//
func (p *OrderList) OnCancelOrder(trackNumber int64, tab string) *exit.Exit {
	return p._onStatusEvent(trackNumber, "canceled", tab)
}

func (p *OrderList) OnDeliver(trackNumber int64, tab string) *exit.Exit {
	return p._onStatusEvent(trackNumber, "delivering", tab)
}

func (p *OrderList) OnMarkAsDone(trackNumber int64, tab string) *exit.Exit {
	return p._onStatusEvent(trackNumber, "done", tab)
}

func (p *OrderList) _onStatusEvent(trackNumber int64, status string, tab string) *exit.Exit {
	err := orderservice.ChangeOrderStatus(trackNumber, status)
	if err != nil {
		panic(err.Error())
	}
	return route.RedirectDispatch(route.GetRefererFromURL(p.Request()), "/order/list")
}

func (p *OrderList) Ondelete(trackNumber int64, tab string) (string, string) {
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
	var editlink string
	switch model.OrderType(order.Type) {
	case model.Wholesale:
		editlink = fmt.Sprintf("/order/create/detail/%v", order.Id)
	case model.ShippingInstead:
		editlink = fmt.Sprintf("/order/create/shippinginstead/%v", order.TrackNumber)
	case model.SubOrder:
		editlink = "#Error"
	default:
		editlink = fmt.Sprintf("%v", order.Type)
	}
	return editlink

	// add return url // another method to return page to add referer
	refer := p.Request().URL.RequestURI()
	return editlink + "?referer=" + refer // TODO: need encode
	// panic(fmt.Sprintf("Wrong order type for %v", order.TrackNumber))
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
		// TODO: Use link generator to generate this link.
		return fmt.Sprintf("/order/list.orderlist:Print/%v", order.TrackNumber)
	case model.ShippingInstead:
		return fmt.Sprintf("/order/list.orderlist:ShippingInsteadOrderPrint/%v", order.TrackNumber)
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
