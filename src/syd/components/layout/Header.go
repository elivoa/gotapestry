package layout

import (
	"fmt"
	"syd/service"

	"github.com/elivoa/got/core"
)

// ________________________________________________________________________________
// Header -- including css and js resources.
//
type Header struct {
	core.Component
	Title  string
	Public bool

	Ng string // enable angularjs libraries.

	// 是使用第一套css还是第二套
	CSS int // css version now support 1 and 2
}

func (c *Header) Setup() {
	if !c.Public {
		// verify user role.
		// 临时加到这里，没登陆是的用户不得查看任何东西。
		fmt.Println("* TRACE USER LOGIN CHECK *************************************************************")
		service.User.RequireRole(c.W, c.R, "admin") // TODO remove w, r. use service injection.
	}
}
