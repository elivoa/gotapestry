package person

import (
	"fmt"
	"syd/dal/accountdao"
	"syd/dal/orderdao"
	"syd/model"
	"syd/service"
	"syd/service/orderservice"
	"time"

	"github.com/elivoa/got/config"
	"github.com/elivoa/got/core"
	"github.com/elivoa/gxl"
)

/* ________________________________________________________________________________
   PersonEdit
*/
type PersonDetail struct {
	core.Page

	Id *gxl.Int `path-param:"1"`

	Current   int `path-param:"2"` // pager: the current item. in pager.
	PageItems int `path-param:"3"` // pager: page size.
	Total     int

	Person      *model.Person
	Orders      []*model.Order
	TodayOrders []*model.Order
}

func (p *PersonDetail) LeavingMessage(order *model.Order) string {
	return orderservice.LeavingMessage(order)
}

func (p *PersonDetail) Setup() {
	if p.Id == nil {
		return
	}
	// fix the pagers
	if p.PageItems <= 0 {
		p.PageItems = config.LIST_PAGE_SIZE // TODO default pager number. Config this.
	}

	var err error
	if p.Person, err = service.Person.GetPersonById(p.Id.Int); err != nil {
		panic(err)
	}
	if p.Person == nil {
		// var err error
		// p.Total, err = orderdao.CountOrderByCustomer("all", p.Person.Id)
		// if err != nil {
		// 	panic(err.Error())
		// }
		// // TODO finish the common conditional query.
		// // query with pager
		// p.Orders, err = orderdao.ListOrderByCustomer(p.Person.Id, "all", p.Current, p.PageItems)
		// if err != nil {
		// 	panic(err.Error())
		// }
	}

	// fetch data
	var parser = service.Order.EntityManager().NewQueryParser()
	parser.Where("customer_id", p.Id.Int)
	parser.Or("type", model.Wholesale, model.ShippingInstead) // restrict type

	// get total
	p.Total, err = parser.Count()
	if err != nil {
		panic(err.Error())
	}

	// 2. get order list.
	parser.Limit(p.Current, p.PageItems) // pager
	p.Orders, err = service.Order.ListOrders(parser, service.WITH_PERSON)
	if err != nil {
		panic(err.Error())
	}
	// Get today orders.
	date := time.Now()
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 0, 1)
	orders, err := orderdao.ListOrderByCustomer_Time(p.Person.Id, start, end)
	if err != nil {
		panic(err.Error())
	}

	orderservice.LoadDetails(orders)
	p.TodayOrders = orders

	// debug print

	fmt.Println("\n\n\n^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	for _, oo := range p.TodayOrders {
		fmt.Println(">>>>>>> order: ", oo.Accumulated, " >>> ", oo)
	}

	// p.TheBigOrder, p.LeavingMessage = orderservice.GenerateLeavingMessage(p.Person.Id, time.Now())
	if true {
		return
	}
}

func (p *PersonDetail) UrlTemplate() string {
	return fmt.Sprintf("/person/detail/%d/{{Start}}/{{PageItems}}", p.Person.Id)
}

func (p *PersonDetail) ShouldShowLeavingMessage(o *model.Order) bool {
	switch model.OrderType(o.Type) {
	// case model.Wholesale, model.ShippingInstead:
	case model.Wholesale, model.SubOrder:
		return true
	}
	return false
}

func (p *PersonDetail) ChangeLogs() []*model.AccountChangeLog {
	changes, err := accountdao.ListAccountChangeLogsByCustomerId(p.Person.Id)
	if err != nil {
		panic(err.Error())
	}
	return changes
}

func (p *PersonDetail) DisplayType(t int) string {
	switch t {
	case 1:
		return "手工修改"
	case 2:
		return "发货"
	case 3:
		return "批量结款"
	case 4:
		return "取消已发货订单，减去累计欠款"
	default:
		return "未知"
	}
}
