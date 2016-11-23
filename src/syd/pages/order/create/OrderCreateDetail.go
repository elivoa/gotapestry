package order

import (
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route"
	"github.com/elivoa/got/route/exit"
	"github.com/elivoa/gxl"
	"syd/model"
	"syd/service"
	"syd/service/orderservice"
	"syd/service/personservice"
)

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

	ReturnThisPage string // form submit return to this page; values: saveonly, ""
}

// init page
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
		// can't edit WholeSale order when it's staus is after delivering.
		if !order.IsStatus("toprint") && order.TypeIs(uint(model.Wholesale)) {
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
	p.Order.Type = p.Type
	p.Order.ParentTrackNumber = p.ParentTrackNumber
	// init person
	p.Customer = personservice.GetCustomer(p.CustomerId)
	if p.Customer == nil {
		panic(fmt.Sprintf("customer not found: id: %v", p.CustomerId))
	}
}

func (p *OrderCreateDetail) IsDaifa() bool {
	return p.Type == 2
}

// before inject submit values, init fields.
func (p *OrderCreateDetail) OnPrepareForSubmit() {
	if p.Id == nil {
		// if create
		p.Order = model.NewOrder()
		p.Order.CustomerId = p.CustomerId
	} else {
		// if edit
		// for security reason, TODO security check here.
		// 读取了数据库的order是为了保证更新的时候不会丢失form中没有的数据；
		o, err := orderservice.GetOrder(p.Id.Int)
		if err != nil {
			panic(err.Error())
		}
		p.Order = o
		// 但是这样做就必须清除form更新的时候需要删除的值，否则form提交和原有值是叠加的，会引起错误；
		// 这里只需要清除列表等数据，这个Order中只有Details是列表。
		p.Order.Details = nil // clear some value
	}
}

// After inject values, do submit.
func (p *OrderCreateDetail) OnSuccess() *exit.Exit {
	// daofu flag // TODO: 这个可以用框架来实现，注入框架实现转换；
	if p.DaoFu == "on" {
		p.Order.ExpressFee = -1
	}

	// fmt.Println(">>> original order details submit from form;")
	// if nil != p.Order.Details {
	// 	for _, d := range p.Order.Details {
	// 		fmt.Println("\t---: ", d.OrderTrackNumber, d.Color, d.Size, " = ", d.Quantity, d.SellingPrice)
	// 	}
	// }
	fmt.Println("\n\n--------------------------------------------------------------------------------")
	if p.ReturnThisPage == "saveonly" {
		// 临时保存，不允许自提情况；
		if p.Order.DeliveryMethod == "TakeAway" {
			p.Order.DeliveryMethod = ""
		}
	}

	var err error
	if p.IsEdit() {
		if p.Order, err = service.Order.UpdateOrder(p.Order); err != nil {
			panic(err)
		}
	} else {
		if p.Order, err = service.Order.CreateOrder(p.Order); err != nil {
			panic(err)
		}
	}

	if p.ReturnThisPage == "saveonly" {
		if p.IsEdit() {
			return nil
		} else {
			// create 需要返回编辑地址
			return exit.Redirect(fmt.Sprintf("/order/create/detail/%d", p.Order.Id))
		}
	} else {
		url := route.GetRefererFromURL(p.R)
		return exit.RedirectFirstValid(url, p.SourceUrl, "/order/list/"+p.Order.Status)
	}
}

func (p *OrderCreateDetail) IsEdit() bool {
	return p.Id != nil
}

func (p *OrderCreateDetail) ProductDetailJson() interface{} {
	return orderservice.OrderDetailsJson(p.Order, false)
}

func (p *OrderCreateDetail) ShowAddress() bool {
	switch model.OrderType(p.Order.Type) {
	case model.SubOrder:
		return true
	default:
		return false
	}
}
