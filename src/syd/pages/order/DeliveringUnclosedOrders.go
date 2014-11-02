package order

import (
	"encoding/json"
	"fmt"
	"github.com/elivoa/got/core"
	"strconv"
	"strings"
	"syd/dal/orderdao"
	"syd/model"
	"syd/service"
	"syd/service/orderservice"
	"time"
)

type DeliveringUnclosedOrders struct {
	core.Page
	CustomerId int `path-param:"1"`
	Orders     []*model.Order
	Referer    string // return here.
}

// default: get all orders of one person
func (p *DeliveringUnclosedOrders) Setup() (string, string) {
	fmt.Print("\n\n\n\n>>>>>>>>>>>>>>>  ;;; referer is , ", p.Referer)
	return ordersJsonByCustomerid(p.CustomerId)
}

// close action
func (p *DeliveringUnclosedOrders) OnbyTrackingNumber(tns string) (string, string) {
	return ordersJsonByTrackNumbers(tns)
	// TODO SYD: strict privileges validation
	// orderservice.BatchCloseOrder(money, customerId)
	// return ordersJsonByTrackNumbers(customerId)
}

// close orders, and return new orderlist with the same parameters.
// TODO::GOT: add t:ac="xxx" to restore activate parameters.
func (p *DeliveringUnclosedOrders) Onbatchclose(money float64, customerId int) (string, string) {
	// TODO SYD: strict privileges validation
	service.Order.BatchCloseOrder(money, customerId)
	return ordersJsonByCustomerid(customerId)
}

// // close orders by trackingnumbers, and return new orderlist with the same parameters.
// func (p *DeliveringUnclosedOrders) Onbatchclose(money float64, trackNumbers string) (string, string) {
// 	// TODO SYD: strict privileges validation
// 	orderservice.BatchCloseOrderByTrackNumbers(money, trackNumbers)
// 	return ordersJson(customerId)
// }

func ordersJsonByTrackNumbers(tns string) (string, string) {
	orders := []*model.Order{}
	pieces := strings.Split(tns, ",")
	for _, piece := range pieces {
		tn, err := strconv.ParseInt(strings.Trim(piece, " "), 10, 64)
		if err != nil {
			panic(err.Error())
		}
		order, err := orderservice.GetOrderByTrackingNumber(tn)
		if err != nil {
			panic(err.Error())
		}
		orders = append(orders, order)
	}
	return toJsonList(orders)
}

func ordersJsonByCustomerid(customerId int) (string, string) {
	orders, err := orderdao.DeliveringUnclosedOrdersByCustomer(customerId)
	if err != nil {
		panic(err.Error())
	}
	return toJsonList(orders)
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

func toJsonList(orders []*model.Order) (string, string) {
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
