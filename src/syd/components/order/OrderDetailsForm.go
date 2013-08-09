package order

import (
	"got/core"
)

type OrderDetailsForm struct {
	core.Component
	HideOperation bool // show operation column?
}
