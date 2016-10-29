package order

import (
	"bytes"
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"html/template"
	"strconv"
	"syd/base"
	"syd/model"
	"syd/service"
	"syd/service/orderservice"
)

// ________________________________________________________________________________
// OrderPrint

type OrderPrint struct {
	core.Page
	TrackNumber int64 `path-param:"1"`

	Order    *model.Order
	Customer *model.Person
	Sumprice float64 // order sum price, no expressfee, no accumulated, no 代发

	//
	First bool `query:"first"`
}

func (p *OrderPrint) Activate() {
	if p.TrackNumber == 0 {
		panic("Need Tracking Number!")
	}
	fmt.Println("Order Print is here")
}

func (p *OrderPrint) Setup() *exit.Exit {
	order, err := orderservice.GetOrderByTrackingNumber(p.TrackNumber)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(">>>>>> 1")
	p.Order = order
	if p.Customer, err = service.Person.GetPersonById(p.Order.CustomerId); err != nil {
		fmt.Println(">>>>>> 2")
		panic(err)
	} else if p.Customer == nil {
		fmt.Println(">>>>>> 3")
		panic("Customer does not exist!")
	}
	fmt.Println(">>>>>> 4")

	if p.First == false {
		fmt.Println(">>>>>> 5")

		needprice := person_need_price(p.Customer.Id)
		if !needprice {
			fmt.Println(">>>>>> 6")

			// Redirect to print no pirce page.
			return exit.Redirect(fmt.Sprintf("/order/printnoprice/%d", p.TrackNumber))
		}
	}
	fmt.Println(">>>>>> 7")

	p.Sumprice = p.sumprice()
	return nil
}

// Return true if the specified person need print price defaultly.
func person_need_price(customerId int) bool {
	// 如果不是跳转过来的，就要检查默认状态。如果用户在不打印价格列表中，就要跳转到不打印文件中。
	result, err := service.Const.Get(base.SYS_PREF_KEY_PRINT_HIDE_PRICE, strconv.Itoa(customerId))
	if err != nil {
		panic(err)
	}
	if nil != result {
		if intvalue, err := result.Get2ndIntValue(); err == nil && intvalue == 0 {
			// need redirect.
			return true
		} else {
			panic(err)
		}
	}
	return false
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

// ________________________________________________________________________________
// Display Summarize
//
func (p *OrderPrint) DeliveryMethodDisplay() string {
	dis, ok := deliveryMethodDisplayMap[p.Order.DeliveryMethod]
	if ok {
		return dis
	} else {
		return p.Order.DeliveryMethod
	}
}

func (p *OrderPrint) DeliveryMethodIs(dm string) bool {
	return p.Order.DeliveryMethod == dm
}

func (p *OrderPrint) HasExpressFee() bool {
	// not 自提 & not 到付， 剩下的就是没填。
	if p.Order.DeliveryMethod != "TakeAway" && p.Order.ExpressFee != -1 {
		return true
	}
	return false
}

func (p *OrderPrint) ExpressFeeHtml() interface{} {
	if p.Order.ExpressFee == 0 {
		return template.HTML("<span class=\"underline\"></span>")
	} else {
		return p.Order.ExpressFee
	}
}

func (p *OrderPrint) TotalPriceHtml() interface{} {
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

func (p *OrderPrint) IsDaofu() bool {
	return p.Order.ExpressFee == -1
}

// not used
func (p *OrderPrint) DeliveryMethodHtml() interface{} {
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
var deliveryMethodDisplayMap = map[string]string{
	"YTO":      "圆通速递",
	"SF":       "顺风快递",
	"Depoon":   "德邦",
	"Freight":  "货运",
	"TakeAway": "自提",
}
