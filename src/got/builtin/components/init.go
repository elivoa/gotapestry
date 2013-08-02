package components

import (
	"got/route"
)

func Register() {
	// used to locate this file's full path in go.
}

func init() {
	route.Component(Register, &Select{}, &ProvinceSelect{})
}
