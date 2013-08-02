package layout

import (
	"got/core"
	"got/route"
)

func Register() {}
func init() {
	route.Component(Register,
		&Header{}, &HeaderNav{}, &LeftNav{},
	)
}

// ________________________________________________________________________________
// Header -- including css and js resources.
//
type Header struct {
	core.Component
	Title string
}

// ________________________________________________________________________________
type HeaderNav struct {
	core.Component
}

// ________________________________________________________________________________
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
