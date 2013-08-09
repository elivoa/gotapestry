/*
  创建/修改订单中的商品选择模块
*/

package order

import (
	"got/core"
)

type OrderProductSelector struct {
	core.Component
	CustomerId int
}

// TODO How to map url to component. like taptestry5
func (p *OrderProductSelector) OnProductJson() string {
	return "product json"
}
