package account

import (
	"github.com/elivoa/got/route/exit"
	"got/core"
	"syd/service"
)

type AccountLogout struct {
	core.Page
}

func (p *AccountLogout) Setup() *exit.Exit {
	service.User.Logout(p.W, p.R) // TODO Need check permission to logout?
	return exit.Redirect("/")
}
