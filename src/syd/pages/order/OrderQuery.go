package order

import (
	"fmt"
	"github.com/elivoa/got/builtin/services"
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"github.com/elivoa/got/utils"
	"strings"
	"syd/model"
	"syd/service"
	"time"
)

/* ________________________________________________________________________________
The Order List page
*/ 
type OrderQuery struct {
	core.Page 

	// parameters
	Orders    []*model.Order
	Tab       string `path-param:"1"`
	Current   int    `path-param:"2"` // pager: the current item. in pager.
	PageItems int    `path-param:"3"` // pager: page size.

	TimeFrom time.Time `query:"from"`
	TimeTo   time.Time `query:"to"`

	// properties
	Total int // pager: total items available
}

func (p *OrderQuery) Activate() {
	// service.User.RequireRole(p.W, p.R, syd.RoleSet_Orders...)

	// not injected with parameters.
	if p.Tab == "" {
		p.Tab = "all" // default go in toprint
	}

	// time
	if !utils.ValidTime(p.TimeTo) && !utils.ValidTime(p.TimeFrom) {
		p.TimeFrom, p.TimeTo = utils.NatureTimeRange(0, 0, -3)
	} else if !utils.ValidTime(p.TimeTo) {
		p.TimeTo = utils.NatureTimeTodayEnd()
	} else if !utils.ValidTime(p.TimeFrom) {
		// TODO what to do?
	}
}

func (p *OrderQuery) SetupRender() {
	// fix the pagers
	if p.PageItems <= 0 {
		p.PageItems = config.LIST_PAGE_SIZE // TODO default pager number. Config this.
	}

	// fetch data
	var err error
	var parser = service.Order.EntityManager().NewQueryParser()
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

	// filter by time range
	if utils.ValidTime(p.TimeTo) && utils.ValidTime(p.TimeFrom) {
		parser.Range("create_time", p.TimeFrom, p.TimeTo)
		// fmt.Println(">>> query by ", p.TimeFrom, p.TimeTo)
	} else if !utils.ValidTime(p.TimeTo) {
	} else if !utils.ValidTime(p.TimeFrom) {
	}

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
}

func (p *OrderQuery) OnSuccessFromSearchForm() *exit.Exit {
	// time is injected and then return linkpage.
	return exit.Redirect(p.ThisPageLink())
}

func (p *OrderQuery) OnClearForm() *exit.Exit {
	p.TimeFrom = time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
	p.TimeTo = time.Now() //  p.TimeFrom
	return exit.Redirect(p.ThisPageLink())
}

func (p *OrderQuery) ThisPageLink() string {
	// 一个普通的SearchBox实现。所有东西都放到url里面。直接redirect到本页面。
	var parameters = map[string]interface{}{
		// "tab":  p.Tab,
		"from": p.TimeFrom,
		"to":   p.TimeTo,
	}
	url := services.Link.GeneratePageUrlWithContextAndQueryParameters("order/query", parameters, p.Tab)
	return url
}

func (p *OrderQuery) OnTab(tab string) *exit.Exit {
	// fmt.Println("============================= p.Tab/tab is ,", p.Tab, "/", tab)
	p.Tab = tab
	return exit.Redirect(p.ThisPageLink())
}

// -------------- not modified ---------------

func (p *OrderQuery) TabStyle(tab string) string {
	if strings.ToLower(p.Tab) == strings.ToLower(tab) {
		return "cur"
	}
	return ""
}

// pager related

func (p *OrderQuery) UrlTemplate() string {
	return fmt.Sprintf("/order/list/%s/{{Start}}/{{PageItems}}", p.Tab)
}
