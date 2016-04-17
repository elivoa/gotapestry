package layout

import (
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/gorilla/sessions"
	"net/http"
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

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func (c *HeaderNav) SetupRender() {
	c.UserToken = service.User.GetLogin(c.W, c.R)
	// Speical Version: don't redirect
	c.TimeZoneInfo = service.TimeZone.UserTimeZoneDontCheckCookie(c.R)

	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	session, err := store.Get(c.R, "session-name")
	if err != nil {
		http.Error(c.W, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Values in session is:")
	for k, v := range session.Values {
		fmt.Println("\tK:", k, " -> ", v)
	}

	// save to session
	session.Values["test"] = "sldkflsdj"
	session.Save(c.R, c.W)
	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")

}

func (c *HeaderNav) IsLogin() bool {
	// c.UserToken = service.User.GetLogin(c.ResponseWriter(), c.Request())
	return c.UserToken != nil
}

func (c *HeaderNav) IsAdmin() bool {
	return c.UserToken.HasRole(base.Role_Admin)
}
