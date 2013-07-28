package order

import (
	"fmt"
	"got/core"
	"got/register"
	"gxl"
	"syd/model"
	"syd/service/orderservice"
	"syd/service/personservice"
)

func init() {
	register.Page(Register, &OrderCreateDetail{})
}

// --------  Order Create Details  -----------------------------------------------------------------------

type OrderCreateDetail struct {
	core.Page
	CustomerId int      `query:"customer"`
	Id         *gxl.Int `path-param:"1"` // OrderId, used when edit.

	Customer *model.Customer
	Order    *model.Order // submit the big table.
	DaoFu    string       // on

	SourceUrl         string `query:"source"` // redirect url
	ParentTrackNumber int64  `query:"parent"` // parent tn if subproject

	// TODO Customized type coercion.
	// Temp Solution: coercion:".CoercionMethod(string)"
	Type uint `query:"type"` // order type

	// page msg resources
	SubTitle     string // create or edit? TODO i18n resource file.
	SubmitButton string // 确认下单？ 修改订单
}

func (p *OrderCreateDetail) Setup() {
	if p.Id == nil && p.CustomerId == 0 {
		panic("Can't find order to edit!")
	}

	// Validate Privileges.
	if p.Id != nil {
		// edit mode
		order, err := orderservice.GetOrder(p.Id.Int)
		if err != nil {
			panic(err.Error())
		}
		if !order.IsStatus("todeliver") {
			panic(fmt.Sprintf("Order are not allow to edit in this status[%v]", order.Status))
		}

		p.Order = order
		p.CustomerId = p.Order.CustomerId
		p.SubTitle = "编辑"
		p.SubmitButton = "修改订单"
	} else {
		// create mode
		p.Order = model.NewOrder()
		p.Order.CustomerId = p.CustomerId
		p.SubTitle = "新建"
		p.SubmitButton = "确认下单"
	}

	// set order type; trick: no set is 0, 0 is default Wholesale type.
	fmt.Println("````````````````````````````````````````````````````````````````````````````````")
	fmt.Println(p.Type)
	p.Order.Type = p.Type
	p.Order.ParentTrackNumber = p.ParentTrackNumber
	// init person
	p.Customer = personservice.GetCustomer(p.CustomerId)
	if p.Customer == nil {
		panic(fmt.Sprintf("customer not found: id: %v", p.CustomerId))
	}
}

// before inject submit values, init fields.
func (p *OrderCreateDetail) OnSubmit() {
	if p.Id == nil {
		// if create
		p.Order = model.NewOrder()
		p.Order.CustomerId = p.CustomerId
	} else {
		// if edit
		// for security reason, TODO security check here.
		o, err := orderservice.GetOrder(p.Id.Int)
		if err != nil {
			panic(err.Error())
		}
		p.Order = o
	}
}

// after submit
func (p *OrderCreateDetail) OnSuccess() (string, string) {
	// order, update details
	for _, detail := range p.Order.Details {
		detail.OrderTrackNumber = p.Order.TrackNumber
	}
	// daofu flag
	if p.DaoFu == "on" {
		p.Order.ExpressFee = -1
	}

	// set new & edited order's status to 'todeliver'
	p.Order.Status = "todeliver"

	// update
	if p.Id != nil {
		orderservice.UpdateOrder(p.Order)
	} else {
		orderservice.CreateOrder(p.Order)
	}

	// return source?
	if p.SourceUrl == "" {
		return "redirect", fmt.Sprintf("/order/print/%v", p.Order.TrackNumber)
	} else {
		return "redirect", p.SourceUrl
	}
}

func (p *OrderCreateDetail) IsEdit() bool {
	return p.Id != nil
}

func (p *OrderCreateDetail) ProductDetailJson() interface{} {
	return orderservice.OrderDetailsJson(p.Order)
}

func (p *OrderCreateDetail) ShowAddress() bool {
	switch model.OrderType(p.Order.Type) {
	case model.SubOrder:
		return true
	default:
		return false
	}
}