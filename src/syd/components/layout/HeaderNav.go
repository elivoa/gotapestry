package layout

import (
	"got/core"
	"syd/model"
	"syd/service"
)

// ________________________________________________________________________________
type HeaderNav struct {
	core.Component
	UserToken *model.UserToken
}

func (c *HeaderNav) IsLogin() bool {
	c.UserToken = service.User.GetLogin(c.ResponseWriter(), c.Request())
	return c.UserToken != nil
}
