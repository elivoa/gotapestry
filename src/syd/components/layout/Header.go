package layout

import (
	"fmt"
	"github.com/elivoa/got/core"
	"syd/service"
)

// ________________________________________________________________________________
// Header -- including css and js resources.
//
type Header struct {
	core.Component
	Title  string
	Public bool

	Ng string // enable angularjs libraries.
}

func (c *Header) Setup() {
	if !c.Public {
		// verify user role.
		// 临时加到这里，没登陆是的用户不得查看任何东西。
		fmt.Println("********************************************************************************")
		service.User.RequireRole(c.W, c.R, "admin") // TODO remove w, r. use service injection.
	}
}
