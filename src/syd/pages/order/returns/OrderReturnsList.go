package returns

import (
	"fmt"
	"strings"
	"syd/model"
	"syd/service"

	"github.com/elivoa/got/config"
	"github.com/elivoa/got/core"
)

/* ________________________________________________________________________________
The Order List page
*/
type OrderReturnsList struct {
	core.Page

	// parameters
	Orders    []*model.Order
	Tab       string `path-param:"1"`
	Current   int    `path-param:"2"` // pager: the current item. in pager.
	PageItems int    `path-param:"3"` // pager: page size.

	// properties
	Total int // pager: total items available
}

func (p *OrderReturnsList) Activate() {
	// service.User.RequireRole(p.W, p.R, syd.RoleSet_Orders...)

	// not injected with parameters.
	if p.Tab == "" {
		p.Tab = "toprint" // default go in toprint
	}
}

func (p *OrderReturnsList) SetupRender() {
	// fix the pagers
	if p.PageItems <= 0 {
		p.PageItems = config.LIST_PAGE_SIZE // TODO default pager number. Config this.
	}

	// fetch data
	var err error
	var parser = service.OrderReturns.EntityManager().NewQueryParser()
	parser.Where()
	switch strings.ToLower(p.Tab) {
	// case "today":
	// 	now := time.Now().UTC()
	// 	start := now.Truncate(time.Hour * 24)
	// 	end := now.AddDate(0, 0, 1).Truncate(time.Hour * 24)
	// 	parser.Where().Range("create_time", start, end)
	// case "returned":
	// 	parser.Where("status", "returned")
	case "all", "":
		// all status
	default:
		parser.And("status", p.Tab)
	}
	parser.Or("type", model.Wholesale, model.ShippingInstead) // restrict type

	// get total
	p.Total, err = parser.Count()
	if err != nil {
		panic(err.Error())
	}

	// 2. get order list.
	parser.Limit(p.Current, p.PageItems) // pager
	p.Orders, err = service.OrderReturns.ListOrders(parser, service.WITH_PERSON)
	if err != nil {
		panic(err.Error())
	}
}

func (p *OrderReturnsList) TabStyle(tab string) string {
	if strings.ToLower(p.Tab) == strings.ToLower(tab) {
		return "cur"
	}
	return ""
}

// pager related

func (p *OrderReturnsList) UrlTemplate() string {
	return fmt.Sprintf("/order/list/%s/{{Start}}/{{PageItems}}", p.Tab)
}
