package order

import (
	"got/core"
	"got/route"
	"syd/model"
)

func init() {
	route.Component(Register, &CustomerInfo{})
}

// --------  Customer Info BLock  -----------------------------------------------------------------------

type CustomerInfo struct {
	core.Component
	Customer *model.Customer
}
