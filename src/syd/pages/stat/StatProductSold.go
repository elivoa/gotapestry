package stat

import (
	"got/core"
)

type StatProductSold struct {
	core.Page
	Days int `path-param:"1" TODO_default:"7"`
}

func (p *StatProductSold) Activate() {
	if p.Days == 0 {
		p.Days = 7
	}
}
