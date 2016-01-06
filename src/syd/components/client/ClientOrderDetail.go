package client

import (
	"fmt"
	"github.com/elivoa/got/core"
	"syd/model"
	"syd/service"
	"time"
)

type ClientOrderDetail struct {
	core.Component

	OrderId int64 // Client Order Id, in inventory table.

	ClientOrder *model.InventoryGroup // order as inventory

	// delete things .
	TotalPrice float64 // all order's price
	Referer    string  // return to this place

	TotalGroups   int
	TotalQuantity int

	// ??? what's this?
	Current   int   `path-param:"2"`   // pager: the current item. in pager.
	PageItems int   `path-param:"3"`   // pager: page size.
	Provider  int64 `query:"provider"` // filter by provider Id.

	// properties
	Total int // pager: total items available

}

func (c *ClientOrderDetail) SetupRender() {
	service.User.RequireRole(c.W, c.R, "admin")
	if c.OrderId <= 0 {
		return
	}

	// 取客户下单订单
	order, err := service.InventoryGroup.GetInventoryGroup(c.OrderId,
		service.WITH_INVENTORIES|service.WITH_PERSON|service.WITH_PRODUCT)
	if err != nil {
		panic(err)
	}
	c.ClientOrder = order

	// 取给客户发货的订单
	var startTime = order.SendTime.Truncate(time.Hour * 24)
	fmt.Println(">> start time is : ", startTime)

	orders, err := getSendingOrders(order)
	fmt.Println(">>> ", orders)
	// calculate total.
	// p.TotalGroups = len(p.InventoryGroups)
	// for _, inv := range p.InventoryGroups {
	// 	if nil != inv {
	// 		p.TotalQuantity += inv.TotalQuantity
	// 	}
	// }
}

func getSendingOrders(order *model.InventoryGroup) ([]*model.Order, error) {

	fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&", order)

	// var err error
	var parser = service.Order.EntityManager().NewQueryParser()
	parser.Where()

	parser.Limit(10) // debug

	// case "today":
	// 	now := time.Now().UTC()
	// 	start := now.Truncate(time.Hour * 24)
	// 	end := now.AddDate(0, 0, 1).Truncate(time.Hour * 24)
	// 	parser.Where().Range("create_time", start, end)
	// case "returned":
	// 	parser.Where("status", "returned")
	// parser.And("status", p.Tab) // any status
	parser.Or("type", model.Wholesale, model.ShippingInstead) // restrict type

	// get total
	// p.Total, err = parser.Count()
	// if err != nil {
	// 	panic(err.Error())
	// }

	// 2. get order list.
	// parser.Limit(p.Current, p.PageItems) // pager
	// p.Orders, err = service.Order.ListOrders(parser, service.WITH_PERSON)
	// if err != nil {
	// 	panic(err.Error())
	// }

	return nil, nil
}

// ________________________________________________________________________________
// Events
//
