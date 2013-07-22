package order

import (
	"encoding/json"
	"fmt"
	"got/core"
	"got/register"
	"syd/dal/orderdao"
	"syd/model"
	"syd/service/orderservice"
	"time"
)

func init() {
	register.Page(Register, &DeliveringUnclosedOrders{})
}

type DeliveringUnclosedOrders struct {
	core.Page
	CustomerId int `path-param:"1"`
	Orders     []*model.Order
}

func (p *DeliveringUnclosedOrders) Setup() (string, string) {
	return ordersJson(p.CustomerId)
}

// TODO::GOT: add t:ac="xxx" to restore activate parameters.
func (p *DeliveringUnclosedOrders) Onbatchclose(money float64, customerId int) (string, string) {
	// TODO SYD: strict privileges validation
	orderservice.BatchCloseOrder(money, customerId)
	return ordersJson(customerId)
}

// ________________________________________________________________________________
// Generate SimpleOrder JSON
//

type SimpleOrders struct {
	Order           []int64
	Orders          map[string]*SimpleOrder
	TotalOrderPrice float64
}

type SimpleOrder struct {
	TrackNumber int64     `json:"tn"`
	CreateTime  time.Time `json:"time"`
	OrderPrice  float64   `json:"price"`
}

func ordersJson(customerId int) (string, string) {
	orders, err := orderdao.DeliveringUnclosedOrdersByCustomer(customerId)
	fmt.Println("-----------------", customerId)
	fmt.Println(orders)
	if err != nil {
		panic(err.Error())
	}
	order := []int64{}
	ordermap := map[string]*SimpleOrder{}
	var totalOrderPrice float64
	for _, o := range orders {
		order = append(order, o.TrackNumber)
		ordermap[fmt.Sprint(o.TrackNumber)] = tosimple(o)
		totalOrderPrice += o.SumOrderPrice()
	}

	sos := &SimpleOrders{
		Order:           order,
		Orders:          ordermap,
		TotalOrderPrice: totalOrderPrice,
	}
	b, err := json.Marshal(sos)
	if err != nil {
		// TODO autodetect: return error json on ajax call.
		return "json", fmt.Sprint(err)
	}
	return "json", string(b)
}

func tosimple(order *model.Order) *SimpleOrder {
	return &SimpleOrder{
		TrackNumber: order.TrackNumber,
		CreateTime:  order.CreateTime,
		OrderPrice:  order.SumOrderPrice(),
	}
}

// ________________________________________________________________________________
// Event::CloseOrder
//
