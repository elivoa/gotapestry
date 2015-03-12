package layout

import (
	"github.com/elivoa/got/core"
)

type LeftNav struct {
	core.Component
	CurPage string
}

func (c *LeftNav) Style(page string) string {
	if page == c.CurPage {
		return "cur"
	}
	return ""
}
