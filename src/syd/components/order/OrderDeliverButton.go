package order

import (
	"got/core"
	"got/register"
)

func init() {
	register.Component(Register, &OrderDeliverButton{})
}

type OrderDeliverButton struct {
	core.Component
	Source string // return to this place

	TrackNumber            int64
	DeliveryMethod         string
	DeliveryTrackingNumber string
	ExpressFee             int64
}
