package order

import (
	"fmt"
	"got/core"
	"got/register"
	"syd/model"
	"syd/service/orderservice"
	"syd/service/personservice"
	"syd/service/productservice"
)

func init() {
	register.Page(Register, &ShippingInstead{})
}

//
// --------  Order Shipping Instead  -----------------------------------------------------------------------
//

type ShippingInstead struct {
	core.Page
	CustomerId  int   `query:"customer"`
	TrackNumber int64 `path-param:"1"` // Used when edit

	Order     *model.Order   // The parent order.
	SubOrders []*model.Order // Sub orders
	Customer  *model.Customer

	SourceUrl    string `query:"source"`   // redirect url
	ReadonlyMode bool   `query:"readonly"` // hide edit button.

	// caches
	productcache map[int]*model.Product
}

func (p *ShippingInstead) Setup() {
	if p.TrackNumber == 0 {
		panic("Can't find order!")
	}
	p.productcache = map[int]*model.Product{}

	// edit mode
	var err error
	if p.Order, err = orderservice.GetOrderByTrackingNumber(p.TrackNumber); err != nil {
		panic(err.Error())
	}
	p.CustomerId = p.Order.CustomerId

	// init person
	p.Customer = personservice.GetCustomer(p.CustomerId)
	if p.Customer == nil {
		panic(fmt.Sprintf("customer not found: id: %v", p.CustomerId))
	}

	// init suborders
	subOrders, err := orderservice.LoadSubOrders(p.Order)
	if err != nil {
		panic(err.Error())
	}
	p.SubOrders = subOrders

	// calculate statistics to parent order
	//   f(x) = Sum(suborder.quantity * unit-price + order.expressfee)
	//
	var totalPrice float64 = 0
	var totalExpressFee int64 = 0
	var totalCount int = 0
	for _, so := range subOrders {
		totalCount += so.TotalCount
		totalPrice += so.TotalPrice //SumOrderPrice()
		if so.ExpressFee > 0 {
			totalExpressFee += so.ExpressFee
		}
	}
	p.Order.TotalCount = totalCount
	p.Order.TotalPrice = totalPrice
	p.Order.ExpressFee = totalExpressFee
}

func (p *ShippingInstead) FirstDetail(order *model.Order) *model.OrderDetail {
	if len(order.Details) > 0 {
		return order.Details[0]
	}
	return nil
}

func (p *ShippingInstead) OtherDetails(order *model.Order) []*model.OrderDetail {
	if len(order.Details) > 1 {
		return order.Details[1:len(order.Details)]
	}
	return nil
}

// func (p *ShippingInstead) OrderTotal
func (p *ShippingInstead) ShowProductName(productId int) string {
	product, ok := p.productcache[productId]
	if ok {
		return product.Name
	} else {
		product := productservice.GetProduct(productId)
		if product != nil {
			p.productcache[productId] = product
			return product.Name
		}
	}
	return fmt.Sprintf("product[%v]", productId)
}

// before submit, here url injection is ready but post data is not
// injected. we get order from db.
func (p *ShippingInstead) OnSubmit() {
	var err error
	if p.Order, err = orderservice.GetOrderByTrackingNumber(p.TrackNumber); err != nil {
		panic(err.Error())
	}
	p.CustomerId = p.Order.CustomerId
}

// after post data is injected, override to p.Order. Thus p.Order here
// is full and ready to persist.
func (p *ShippingInstead) OnSuccess() (string, string) {
	orderservice.UpdateOrder(p.Order)
	return "redirect", p.ThisPage()
}

// --------  Links  ------------------------------------------------------------------------

func (p *ShippingInstead) CreateSubOrderLink() string {
	return fmt.Sprintf("/order/create/detail?customer=%v&source=%v&parent=%v&type=%v",
		p.CustomerId, p.ThisPage(), p.TrackNumber, model.SubOrder)
}

func (p *ShippingInstead) EditLink(order *model.Order) string {
	return fmt.Sprintf("/order/create/detail/%v?source=%v&parent=%v&type=%v",
		order.Id, p.ThisPage(), p.TrackNumber, model.SubOrder)
}

func (p *ShippingInstead) ThisPage() string {
	return fmt.Sprintf("/order/create/shippinginstead/%v", p.TrackNumber)
}
