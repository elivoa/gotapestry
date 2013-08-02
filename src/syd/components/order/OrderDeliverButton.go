package order

import (
	"got/core"
	"got/route"
)

func init() {
	route.Component(Register, &OrderDeliverButton{})
}

type OrderDeliverButton struct {
	core.Component
	Source string // return to this place

	TrackNumber            int64
	DeliveryMethod         string
	DeliveryTrackingNumber string
	ExpressFee             int64
}
