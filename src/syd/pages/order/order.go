package order

import (
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route"
	"github.com/elivoa/got/route/exit"
	"strings"
	"syd/dal/accountdao"
	"syd/model"
	"syd/service"
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
	TrackNumber int64 // our order tracknumber

	// need by deliver from
	DeliveryMethod         string
	DeliveryTrackingNumber string
	ExpressFee             int64
	DaoFu                  string

	// need by close form
	Money   float64
	Referer string
}

// **** important logic ****
// TODO transaction. Move to right place. 发货
func (p *ButtonSubmitHere) OnSuccessFromDeliverForm() *exit.Exit {

	var expressFee int64 = 0
	if p.DaoFu == "on" {
		// if order.ExpressFee == -1, means this is `daofu`, don't add -1 to 累计欠款.
		// TODO add field isDaofu to order table. Change ExpressFee to 0;
		expressFee = -1
	} else {
		expressFee = p.ExpressFee
	}

	if _, err := service.Order.DeliverOrder(
		p.TrackNumber, p.DeliveryTrackingNumber, p.DeliveryMethod, expressFee); err != nil {
		panic(err)
	}
	return route.RedirectDispatch(p.Referer, "/order/list")

	if false { // backup, has been replace with above.

		/////////////

		// 1/2 update delivery informantion to order.

		// 1. get order form db.
		order, err := orderservice.GetOrderByTrackingNumber(p.TrackNumber)
		if err != nil {
			panic(err.Error())
		}

		// 2. set data back to order.
		order.DeliveryTrackingNumber = p.DeliveryTrackingNumber
		order.DeliveryMethod = p.DeliveryMethod
		if p.DaoFu == "on" {
			// if order.ExpressFee == -1, means this is `daofu`, don't add -1 to 累计欠款.
			// TODO add field isDaofu to order table. Change ExpressFee to 0;
			order.ExpressFee = -1
		} else {
			order.ExpressFee = p.ExpressFee
		}
		order.Status = "delivering"

		// 3. get person, check if customer exists.
		customer, err := service.Person.GetPersonById(order.CustomerId)
		if err != nil {
			panic(err)
		} else if customer == nil {
			panic(fmt.Sprintf("Customer not found for order! id %v", order.CustomerId))
		}

		// 4. the last chance to update accumulated.
		order.Accumulated = -customer.AccountBallance

		// 5. save order changes.
		if _, err := service.Order.UpdateOrder(order); err != nil {
			panic(err.Error())
		}

		// 6. update customer's AccountBallance
		switch model.OrderType(order.Type) {
		case model.Wholesale, model.SubOrder: // 代发不参与, 代发订单由其子订单负责参与累计欠款的统计；
			customer.AccountBallance -= order.TotalPrice
			if order.ExpressFee > 0 {
				customer.AccountBallance -= float64(order.ExpressFee)
			}
			if _, err = personservice.Update(customer); err != nil {
				panic(err.Error())
			}

			// create chagne log.
			accountdao.CreateAccountChangeLog(&model.AccountChangeLog{
				CustomerId:     customer.Id,
				Delta:          -order.TotalPrice,
				Account:        customer.AccountBallance,
				Type:           2, // order.send
				RelatedOrderTN: order.TrackNumber,
				Reason:         "",
			})

		}
		fmt.Println(">>>>>>>>>>>>>>>>>>>> update all done......................")
	}
	return nil
}

// **** important logic ****
// when close order. 结款， Close Order
func (p *ButtonSubmitHere) OnSuccessFromCloseForm() *exit.Exit {
	fmt.Println("\n\n\n>>>>>>>>>>>>>>>>>>>> On success from close order.................", p.TrackNumber)
	fmt.Println(">>> ", p.Referer, "<<<<<<<<<<<<<<<<<")

	// 1/2 update delivery informantion to order.
	order, err := orderservice.GetOrderByTrackingNumber(p.TrackNumber)
	if err != nil {
		panic(err.Error())
	}
	order.Status = "done"
	_, err = service.Order.UpdateOrder(order)
	if err != nil {
		panic(err.Error())
	}

	// 2/2 update customer's AccountBallance
	customer, err := service.Person.GetPersonById(order.CustomerId)
	if err != nil {
		panic(err)
	}
	if customer == nil {
		panic(fmt.Sprintf("Customer not found for order! id %v", order.CustomerId))
	}
	customer.AccountBallance += p.Money
	if _, err = personservice.Update(customer); err != nil {
		panic(err.Error())
	}

	// create chagne log at the same time:
	accountdao.CreateAccountChangeLog(&model.AccountChangeLog{
		CustomerId:     customer.Id,
		Delta:          p.Money,
		Account:        customer.AccountBallance,
		Type:           3, // batch close order.
		RelatedOrderTN: order.TrackNumber,
		Reason:         "",
	})
	return route.RedirectDispatch(p.Referer, "/order/list")
}
