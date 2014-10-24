package order

import (
	"github.com/elivoa/got/core"
)

type OrderDeliverButton struct {
	core.Component
	Source string // return to this place

	TrackNumber            int64
	DeliveryMethod         string
	DeliveryTrackingNumber string
	ExpressFee             int64
}
