package components

import (
	"got/register"
)

func Register() {
	// used to locate this file's full path in go.
}

func init() {
	register.Component(Register, &Select{}, &ProvinceSelect{})
}
