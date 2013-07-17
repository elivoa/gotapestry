package order

import (
	"fmt"
	"got/core"
	"got/register"
	"syd/model"
	"syd/service/orderservice"
	"syd/service/personservice"
)

func init() {
	register.Page(Register, &OrderPrint{})
}

// ________________________________________________________________________________
// OrderPrint

type OrderPrint struct {
	core.Page
	TrackNumber string `path-param:"1"`

	Order    *model.Order
	Customer *model.Person
	Sumprice float64
}

func (p *OrderPrint) Activate() {
	if p.TrackNumber == "" {
		panic("Need Tracking Number!")
	}
	fmt.Println("Order Print is here")
}

func (p *OrderPrint) Setup() {
	order, err := orderservice.GetOrderByTrackingNumber(p.TrackNumber)
	if err != nil {
		panic(err.Error())
	}
	p.Order = order
	if p.Customer = personservice.GetPerson(p.Order.CustomerId); p.Customer == nil {
		panic("Customer does not exist!")
	}
	p.Sumprice = p.sumprice()
}

func (p *OrderPrint) DeliveryMethodDisplay() string {
	dis, ok := deliveryMethodDisplayMap[p.Order.DeliveryMethod]
	if ok {
		return dis
	} else {
		return p.Order.DeliveryMethod
	}
}

func (p *OrderPrint) sumprice() float64 {
	var sum float64
	if p.Order.Details != nil {
		for _, detail := range p.Order.Details {
			sum = sum + float64(detail.Quantity)*detail.SellingPrice
		}
	}
	return sum
}
func (p *OrderPrint) ProductDetailJson() interface{} {
	return orderservice.OrderDetailsJson(p.Order)
}

var deliveryMethodDisplayMap = map[string]string{
	"YTO":      "圆通速递",
	"SF":       "顺风快递",
	"TakeAway": "自提",
}
