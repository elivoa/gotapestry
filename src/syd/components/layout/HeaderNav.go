package layout

import (
	"github.com/elivoa/got/core"
	"syd/base"
	"syd/model"
	"syd/service"
)

// ________________________________________________________________________________
type HeaderNav struct {
	core.Component
	UserToken    *model.UserToken
	TimeZoneInfo *model.TimeZoneInfo
}

func (c *HeaderNav) SetupRender() {
	c.UserToken = service.User.GetLogin(c.W, c.R)
	// Speical Version: don't redirect
	c.TimeZoneInfo = service.TimeZone.UserTimeZoneDontCheckCookie(c.R)
}

func (c *HeaderNav) IsLogin() bool {
	// c.UserToken = service.User.GetLogin(c.ResponseWriter(), c.Request())
	return c.UserToken != nil
}

func (c *HeaderNav) IsAdmin() bool {
	return c.UserToken.HasRole(base.Role_Admin)
}
