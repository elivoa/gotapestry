package stat

import (
	"github.com/elivoa/got/core"
)

type StatProductSold struct {
	core.Page
	Days int `path-param:"1" TODO_default:"7"`
}

func (p *StatProductSold) Activate() {
	if !p.Injected("Days") {
		p.Days = 7
	}
}
