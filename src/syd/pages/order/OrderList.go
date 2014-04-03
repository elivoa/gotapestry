package order

import (
	"fmt"
	"got/core"
	"strings"
	"syd/model"
	"syd/service/orderservice"
)

/* ________________________________________________________________________________
The Order List page
*/
type OrderList struct {
	core.Page

	// parameters
	Orders    []*model.Order
	Tab       string `path-param:"1"`
	Current   int    `path-param:"2"` // pager: the current item. in pager.
	PageItems int    `path-param:"3"` // pager: page size.

	// properties
	Total int // pager: total items available

	// customerNames map[int]*model.Person // order-id -> customer names
}

func (p *OrderList) Activate() {
	// not injected with parameters.

	// fix parameters
	if p.Tab == "" {
		p.Tab = "toprint" // default go in toprint
	}
}

func (p *OrderList) SetupRender() {
	// fix the pagers
	if p.PageItems <= 0 {
		p.PageItems = 50 // TODO default pager number. Config this.
	}

	// 1. get total
	// 2. get order list.
	var err error
	p.Total, err = orderservice.CountOrder(p.Tab)
	if err != nil {
		panic(err.Error())
	}

	p.Orders, err = orderservice.ListOrderPager(p.Tab, p.Current, p.PageItems)
	if err != nil {
		panic(err.Error())
	}

	// p.Orders = dal.ListOrder(p.Tab)
}

func (p *OrderList) TabStyle(tab string) string {
	if strings.ToLower(p.Tab) == strings.ToLower(tab) {
		return "cur"
	}
	return ""
}

// pager related

func (p *OrderList) UrlTemplate() string {
	return fmt.Sprintf("/order/list/%s/{{Start}}/{{PageItems}}", p.Tab)
}
