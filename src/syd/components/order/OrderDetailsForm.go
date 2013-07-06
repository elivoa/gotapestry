package order

import (
	"got/core"
	"got/register"
)

func init() {
	register.Component(Register, &OrderDetailsForm{})
}

type OrderDetailsForm struct {
	core.Component
}
