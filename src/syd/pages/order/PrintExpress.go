package order

import (
	"got/core"
	"got/register"
)

func init() {
	register.Page(Register, &PrintExpressYTO{}, &PrintExpressSF{})
}

// TODO how to solve this problem, 1 template with two different templates.
// TODO!!: Doesn't support to inject to composed component.

// this is base
type printExpress struct {
	core.Page
	DeliveryMethod string `path-param:"1"`
	Address        string `query:"address"`
}

func (p *printExpress) Setup() {

}

// --------------------------------------------------------------------------------
type PrintExpressYTO struct {
	printExpress
	Address string `query:"address"`
}

// --------------------------------------------------------------------------------
type PrintExpressSF struct {
	printExpress
	Address string `query:"address"`
}
