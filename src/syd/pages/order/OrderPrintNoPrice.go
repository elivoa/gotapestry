package order

import (
	"bytes"
	"fmt"
	"github.com/elivoa/got/core"
	"html/template"
	"syd/model"
	"syd/service"
	"syd/service/orderservice"
)

// ________________________________________________________________________________
// OrderPrint

type OrderPrintNoPrice struct {
	core.Page
	TrackNumber int64 `path-param:"1"`

	Order    *model.Order
	Customer *model.Person
	Sumprice float64 // order sum price, no expressfee, no accumulated, no 代发
}

func (p *OrderPrintNoPrice) Activate() {
	if p.TrackNumber == 0 {
		panic("Need Tracking Number!")
	}
	fmt.Println("Order Print is here")
}

func (p *OrderPrintNoPrice) Setup() {
	order, err := orderservice.GetOrderByTrackingNumber(p.TrackNumber)
	if err != nil {
		panic(err.Error())
	}

	p.Order = order
	if p.Customer, err = service.Person.GetPersonById(p.Order.CustomerId); err != nil {
		panic(err)
	} else if p.Customer == nil {
		panic("Customer does not exist!")
	}

	p.Sumprice = p.sumprice()
}

func (p *OrderPrintNoPrice) sumprice() float64 {
	var sum float64
	if p.Order.Details != nil {
		for _, detail := range p.Order.Details {
			sum = sum + float64(detail.Quantity)*detail.SellingPrice
		}
	}
	return sum
}

func (p *OrderPrintNoPrice) ProductDetailJson() interface{} {
	return orderservice.OrderDetailsJson(p.Order)
}

// ________________________________________________________________________________
// Display Summarize
//
func (p *OrderPrintNoPrice) DeliveryMethodDisplay() string {
	dis, ok := deliveryMethodDisplayMap[p.Order.DeliveryMethod]
	if ok {
		return dis
	} else {
		return p.Order.DeliveryMethod
	}
}

func (p *OrderPrintNoPrice) DeliveryMethodIs(dm string) bool {
	return p.Order.DeliveryMethod == dm
}

func (p *OrderPrintNoPrice) HasExpressFee() bool {
	// not 自提 & not 到付， 剩下的就是没填。
	if p.Order.DeliveryMethod != "TakeAway" && p.Order.ExpressFee != -1 {
		return true
	}
	return false
}

func (p *OrderPrintNoPrice) ExpressFeeHtml() interface{} {
	if p.Order.ExpressFee == 0 {
		return template.HTML("<span class=\"underline\"></span>")
	} else {
		return p.Order.ExpressFee
	}
}

func (p *OrderPrintNoPrice) TotalPriceHtml() interface{} {
	// 自提 到付， 显示总订单额就好
	if p.Order.DeliveryMethod != "TakeAway" && p.Order.ExpressFee != -1 {
		return p.Sumprice
	}
	if p.Order.ExpressFee == 0 { // 没填
		return template.HTML("<span class=\"underline\"></span>")
	} else {
		return float64(p.Order.ExpressFee) + p.Sumprice // 合计
	}
}

func (p *OrderPrintNoPrice) IsDaofu() bool {
	return p.Order.ExpressFee == -1
}

// not used
func (p *OrderPrintNoPrice) DeliveryMethodHtml() interface{} {
	var html bytes.Buffer
	html.WriteString(p.DeliveryMethodDisplay())
	html.WriteString("        ")
	if p.Order.ExpressFee == -1 {
		html.WriteString("到付")
	} else {
		html.WriteString("运费: ")
		if p.Order.ExpressFee == 0 {
			html.WriteString("<span class=\"underline\"></span>")
			// html.WriteString("______________")
		} else {
			html.WriteString(fmt.Sprintf("%v", p.Order.ExpressFee))
		}
	}
	return template.HTML(html.String())
}
