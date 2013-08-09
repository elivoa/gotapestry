package order

import (
	"got/core"
	"syd/model"
)

type CustomerInfo struct {
	core.Component
	Customer *model.Customer
}
