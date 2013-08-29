package order

import (
	"got/core"
	"time"
)

// TODO how to solve this problem, 1 template with two different templates.
// TODO!!: Doesn't support to inject to composed component.
// TODO: Maybe a Bug: Pages Extend, parser can't find class that indirectly implement core.Page.

// this is base
type printExpress struct {
	core.Page
	DeliveryMethod string `path-param:"1"`
	Address        string `query:"address"`
	Sender         string `query:"sender"`
}

var senderAddress = map[string]string{
	"李玉勋": "李玉勋，13004211905",
	"王爽":  "王爽，18638068666",
	"小郁":  "mimi，13918116067",
	"吴凤仙": "薛芳芝，13918613475",
	"黄继华": "钱国英，13567338479",
}

// --------------------------------------------------------------------------------
type PrintExpressYTO struct {
	// printExpress // TODO: injection: can't find embed struct to set value.
	core.Page
	DeliveryMethod string `path-param:"1"`
	Address        string `query:"address"`
	Sender         string `query:"sender"`
	Quantity       int    `query:"quantity"`
}

func (p *PrintExpressYTO) SenderAddress() string {
	address := senderAddress[p.Sender]
	if address == "" {
		return " ， "
	}
	return address
}

func (p *PrintExpressYTO) Today() time.Time {
	return time.Now()
}

// --------------------------------------------------------------------------------

type PrintExpressSF struct {
	printExpress
	Address string `query:"address"`
}
