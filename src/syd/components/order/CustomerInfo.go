package order

import (
	"got/core"
	"got/register"
	"syd/model"
)

func init() {
	register.Component(Register, &CustomerInfo{})
}

// --------  Customer Info BLock  -----------------------------------------------------------------------

type CustomerInfo struct {
	core.Component
	Customer *model.Customer
}
