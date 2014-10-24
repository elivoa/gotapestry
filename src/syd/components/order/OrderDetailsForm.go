package order

import (
	"github.com/elivoa/got/core"
)

type OrderDetailsForm struct {
	core.Component
	HideOperation bool // show operation column?
}
