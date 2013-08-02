package order

import (
	"got/core"
	"got/route"
)

func init() {
	route.Component(Register, &OrderDetailsForm{})
}

type OrderDetailsForm struct {
	core.Component
	HideOperation bool // show operation column?
}
