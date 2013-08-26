package order

import (
	"fmt"
	"got/core"
	"strings"
	"syd/model"
	"syd/service/orderservice"
	"syd/service/personservice"
)

// ________________________________________________________________________________
// Start Order pages
//
// ____ Order Index _______________________________________________________________
type OrderIndex struct {
	core.Page           `PageRedirect:"/order/list"`
	__got_page_redirect int `PageRedirect:"/order/list"` //?
}

func (p *OrderIndex) SetupRender() (string, string) {
	return "redirect", "/order/list"
}

/* ________________________________________________________________________________
   Order List
*/
type OrderList struct {
	core.Page

	Orders []*model.Order
	Tab    string `path-param:"1"`

	// customerNames map[int]*model.Person // order-id -> customer names
}

func (p *OrderList) Activate() {
	if p.Tab == "" {
		p.Tab = "toprint" // default go in toprint
	}
}

func (p *OrderList) SetupRender() {
	orders, err := orderservice.ListOrder(p.Tab)
	if err != nil {
		panic(err.Error())
	}
	p.Orders = orders
	// p.Orders = dal.ListOrder(p.Tab)
}

func (p *OrderList) TabStyle(tab string) string {
	if strings.ToLower(p.Tab) == strings.ToLower(tab) {
		return "cur"
	}
	return ""
}

// EVENT: cancel order.
// TODO: put this on component.
// TODO: return null to refresh the current page.
func (p *OrderList) OnCancelOrder(trackNumber int64, tab string) (string, string) {
	return p._onStatusEvent(trackNumber, "canceled", tab)
}

func (p *OrderList) OnDeliver(trackNumber int64, tab string) (string, string) {
	return p._onStatusEvent(trackNumber, "delivering", tab)
}

func (p *OrderList) OnMarkAsDone(trackNumber int64, tab string) (string, string) {
	return p._onStatusEvent(trackNumber, "done", tab)
}

func (p *OrderList) _onStatusEvent(trackNumber int64, status string, tab string) (string, string) {
	err := orderservice.ChangeOrderStatus(trackNumber, status)
	if err != nil {
		panic(err.Error())
	}
	return "redirect", "/order/list/" + tab
}

/* ________________________________________________________________________________
   Order List
*/
type ShippingInsteadList struct {
	core.Page

	Orders []*model.Order
	Tab    string `path-param:"1"`

	// customerNames map[int]*model.Person // order-id -> customer names
}

func (p *ShippingInsteadList) Activate() {
	if p.Tab == "" {
		p.Tab = "todeliver" // default go in todeliver
	}
}

func (p *ShippingInsteadList) SetupRender() {
	orders, err := orderservice.ListOrderByType(model.ShippingInstead, p.Tab)
	if err != nil {
		panic(err.Error())
	}
	p.Orders = orders
}

func (p *ShippingInsteadList) TabStyle(tab string) string {
	if strings.ToLower(p.Tab) == strings.ToLower(tab) {
		return "cur"
	}
	return ""
}

// EVENT: cancel order.
// TODO: put this on component.
// TODO: return null to refresh the current page.
func (p *ShippingInsteadList) OnCancelOrder(trackNumber int64, tab string) (string, string) {
	return p._onStatusEvent(trackNumber, "canceled", tab)
}

func (p *ShippingInsteadList) OnDeliver(trackNumber int64, tab string) (string, string) {
	return p._onStatusEvent(trackNumber, "delivering", tab)
}

func (p *ShippingInsteadList) OnMarkAsDone(trackNumber int64, tab string) (string, string) {
	return p._onStatusEvent(trackNumber, "done", tab)
}

func (p *ShippingInsteadList) _onStatusEvent(trackNumber int64, status string, tab string) (string, string) {
	err := orderservice.ChangeOrderStatus(trackNumber, status)
	if err != nil {
		panic(err.Error())
	}
	return "redirect", "/order/list/" + tab
}

// ________________________________________________________________________________
// ________________________________________________________________________________
// EVNET: Form submits here
// TODO: 功能按钮的表单暂时提交到这里，因为组件内提交暂时还没做好。TODO 快吧组件功能实现了吧。
//
type ButtonSubmitHere struct {
	core.Page
	Source      string
	TrackNumber int64 // our order tracknumber

	// need by deliver from
	DeliveryMethod         string
	DeliveryTrackingNumber string
	ExpressFee             int64
	DaoFu                  string

	// need by close form
	Money float64
}

// **** important logic ****
// TODO transaction. Move to right place
func (p *ButtonSubmitHere) OnSuccessFromDeliverForm() (string, string) {
	// 1/2 update delivery informantion to order.

	fmt.Println(">>>>>>>>>>>>>>>>>>>> update order......................")
	order, err := orderservice.GetOrderByTrackingNumber(p.TrackNumber)
	if err != nil {
		panic(err.Error())
	}
	order.DeliveryTrackingNumber = p.DeliveryTrackingNumber
	order.DeliveryMethod = p.DeliveryMethod
	if p.DaoFu == "on" {
		order.ExpressFee = -1
	} else {
		order.ExpressFee = p.ExpressFee
	}
	order.Status = "delivering"
	_, err = orderservice.UpdateOrder(order)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(">>>>>>>>>>>>>>>>>>>> update pesron......................")

	// 2/2 update customer's AccountBallance
	switch model.OrderType(order.Type) {
	case model.Wholesale, model.SubOrder: //
		customer := personservice.GetPerson(order.CustomerId)
		if customer == nil {
			panic(fmt.Sprintf("Customer not found for order! id %v", order.CustomerId))
		}
		customer.AccountBallance -= order.TotalPrice
		if order.ExpressFee > 0 {
			customer.AccountBallance -= float64(order.ExpressFee)
		}
		if _, err = personservice.Update(customer); err != nil {
			panic(err.Error())
		}
	}
	fmt.Println(">>>>>>>>>>>>>>>>>>>> update all done......................")

	fmt.Println("_____ on deliver from devliver from --- success --- ")
	return p.returnDispatch()
}

// **** important logic ****
func (p *ButtonSubmitHere) OnSuccessFromCloseForm() (string, string) {
	// 1/2 update delivery informantion to order.
	order, err := orderservice.GetOrderByTrackingNumber(p.TrackNumber)
	if err != nil {
		panic(err.Error())
	}
	order.Status = "done"
	_, err = orderservice.UpdateOrder(order)
	if err != nil {
		panic(err.Error())
	}

	// 2/2 update customer's AccountBallance
	customer := personservice.GetPerson(order.CustomerId)
	if customer == nil {
		panic(fmt.Sprintf("Customer not found for order! id %v", order.CustomerId))
	}
	customer.AccountBallance += p.Money
	if _, err = personservice.Update(customer); err != nil {
		panic(err.Error())
	}

	return p.returnDispatch()
}

func (p *ButtonSubmitHere) returnDispatch() (string, string) {
	if p.Source != "" {
		return "redirect", p.Source
	} else {
		return "redirect", "/order/list/all"
	}

}

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
	p.Customer = personservice.GetPerson(p.Order.CustomerId)
	if p.Customer == nil {
		panic(fmt.Sprintf("customer not found: id: %v", p.Order.CustomerId))
	}
}

func (p *ViewOrder) ProductDetailJson() interface{} {
	return orderservice.OrderDetailsJson(p.Order)
}
