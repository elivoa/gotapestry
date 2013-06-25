package layout

import (
	"got/core"
	"got/register"
)

func Register() {
	register.Component(Register,
		&Header{}, &LeftNav{},
	)
}

type Header struct {
	core.Component
}

type LeftNav struct {
	core.Component
}
