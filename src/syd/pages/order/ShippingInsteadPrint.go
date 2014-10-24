package order

import (
	"bytes"
	"fmt"
	"github.com/elivoa/got/core"
	"html/template"
	"syd/model"
	"syd/service/orderservice"
	"syd/service/personservice"
	"syd/service/productservice"
)

// ________________________________________________________________________________
// ShippingInsteadPrint

type ShippingInsteadPrint struct {
	core.Page
	TrackNumber int64 `path-param:"1"`

	Order     *model.Order
	SubOrders []*model.Order // Sub orders
	Customer  *model.Person

	// caches
	productcache map[int]*model.Product

	Sumprice float64 // order sum price, no expressfee, no accumulated, no 代发
}

func (p *ShippingInsteadPrint) Activate() {
	if p.TrackNumber == 0 {
		panic("Need Tracking Number!")
	}
	fmt.Println("Order Print is here")
}

func (p *ShippingInsteadPrint) Setup() {
	p.productcache = map[int]*model.Product{}

	order, err := orderservice.GetOrderByTrackingNumber(p.TrackNumber)
	if err != nil {
		panic(err.Error())
	}
	p.Order = order
	if p.Customer = personservice.GetPerson(p.Order.CustomerId); p.Customer == nil {
		panic("Customer does not exist!")
	}

	// killthis ?
	// logic: update order's accumulated
	if p.Order.Accumulated != -p.Customer.AccountBallance {
		p.Order.Accumulated = -p.Customer.AccountBallance
		_, err := orderservice.UpdateOrder(p.Order)
		if err != nil {
			panic(err.Error())
		}
	}
	// kill this?
	p.Sumprice = p.sumprice()

	// init suborders
	subOrders, err := orderservice.LoadSubOrders(p.Order)
	if err != nil {
		panic(err.Error())
	}
	p.SubOrders = subOrders

	// TODO use calculate or use value in db?
	// calculate statistics to parent order
	//   f(x) = Sum(suborder.quantity * unit-price + order.expressfee)
	//
	var totalPrice float64 = 0
	var totalExpressFee int64 = 0
	var totalCount int = 0
	for _, so := range subOrders {
		totalCount += so.TotalCount
		totalPrice += so.SumOrderPrice()
		if so.ExpressFee > 0 {
			totalExpressFee += so.ExpressFee
		}
	}
	p.Order.TotalCount = totalCount
	p.Order.TotalPrice = totalPrice
	p.Order.ExpressFee = totalExpressFee

}

// func (p *ShippingInstead) OrderTotal
func (p *ShippingInsteadPrint) ShowProductName(productId int) string {
	product, ok := p.productcache[productId]
	if ok {
		return product.Name
	} else {
		product := productservice.GetProduct(productId)
		if product != nil {
			p.productcache[productId] = product
			return product.Name
		}
	}
	return fmt.Sprintf("product[%v]", productId)
}

func (p *ShippingInsteadPrint) sumprice() float64 {
	var sum float64
	if p.Order.Details != nil {
		for _, detail := range p.Order.Details {
			sum = sum + float64(detail.Quantity)*detail.SellingPrice
		}
	}
	return sum
}

func (p *ShippingInsteadPrint) ProductDetailJson() interface{} {
	return orderservice.OrderDetailsJson(p.Order)
}

// ________________________________________________________________________________
// Display Summarize
//
func (p *ShippingInsteadPrint) DeliveryMethodDisplay() string {
	dis, ok := deliveryMethodDisplayMap[p.Order.DeliveryMethod]
	if ok {
		return dis
	} else {
		return p.Order.DeliveryMethod
	}
}

func (p *ShippingInsteadPrint) DeliveryMethodIs(dm string) bool {
	return p.Order.DeliveryMethod == dm
}

func (p *ShippingInsteadPrint) HasExpressFee() bool {
	// not 自提 & not 到付， 剩下的就是没填。
	if p.Order.DeliveryMethod != "TakeAway" && p.Order.ExpressFee != -1 {
		return true
	}
	return false
}

func (p *ShippingInsteadPrint) ExpressFeeHtml() interface{} {
	if p.Order.ExpressFee == 0 {
		return template.HTML("<span class=\"underline\"></span>")
	} else {
		return p.Order.ExpressFee
	}
}

func (p *ShippingInsteadPrint) TotalPriceHtml() interface{} {
	// 自提 到付， 显示总订单额就好
	if p.Order.DeliveryMethod != "TakeAway" && p.Order.ExpressFee != -1 {
		return p.Sumprice
	}
	if p.Order.ExpressFee == 0 { // 没填
		return template.HTML("<span class=\"underline\"></span>")
	} else {
		if p.Order.ExpressFee > 0 {
			return float64(p.Order.ExpressFee) + p.Sumprice // 合计
		} else {
			return "~ERROR~"
		}
	}
}

func (p *ShippingInsteadPrint) IsDaofu() bool {
	return p.Order.ExpressFee == -1
}

// not used
func (p *ShippingInsteadPrint) DeliveryMethodHtml() interface{} {
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

// ________________________________________________________________________________
// helper
//
// in print order
// var deliveryMethodDisplayMap = map[string]string{
// 	"YTO":      "圆通速递",
// 	"SF":       "顺风快递",
// 	"TakeAway": "自提",
// }
