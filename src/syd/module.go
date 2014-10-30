package syd

import (
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/coreservice/sessions"
	"github.com/elivoa/got/errorhandler"
	"github.com/elivoa/got/templates"
	"github.com/elivoa/got/utils"
	"net/http"
	"reflect"
	"strings"
	"syd/model"
)

// todo: think out a better way to register this.
var SYDModule = &core.Module{
	Name:            "syd",       // Don't use this. It's only used to display.
	Version:         "3.0",       // TODO: used to add to assets path to disable cache.
	VarName:         "SYDModule", // Variable name.
	BasePath:        utils.CurrentBasePath(),
	PackagePath:     "syd", // package name used anywhere to locate important things.
	Description:     "SYD Selling System Main module.",
	IsStartupModule: true,
	Register: func() {
		c := config.Config

		// config static resources
		c.AddStaticResource("/pictures/", "/var/site/data/syd/pictures/")
		c.AddStaticResource("/static/", "../static/") // TODO: test this, is this works now?

		c.Port = 8080
		c.DBPort = 3306
		c.DBName = "syd"
		c.DBUser = "root"
		c.DBPassword = "eserver409$)("

		// builtin functions
		templates.RegisterFunc("HasAnyRole", HasAnyRole)

		// errorhandlers
		errorhandler.AddHandler("LoginError",
			reflect.TypeOf(LoginError{}),
			errorhandler.RedirectHandler("/account/login"),
		)

		// --------------------------------------------------------------------------------
		config.LIST_PAGE_SIZE = 50

	},
}

func HasAnyRole(w http.ResponseWriter, r *http.Request, roles ...string) bool {
	session := sessions.LongCookieSession(r)
	if userTokenRaw, ok := session.Values[config.USER_TOKEN_SESSION_KEY]; ok && userTokenRaw != nil {
		if userToken := userTokenRaw.(*model.UserToken); userToken != nil {
			// TODO: check if userToken is outdated.
			if outdated := false; !outdated {
				// TODO: update userToken.Tiemout
				// userToken := service.UserService.GetLogin(w, r)
				if userToken.Roles != nil {
					for _, requiredRole := range roles {
						requiredRole = strings.ToLower(requiredRole)
						for _, role := range userToken.Roles {
							if strings.ToLower(role) == requiredRole {
								return true
							}
						}
					}
				}

			}
		}
	}
	return false
}

// permisstion denied error
type AccessDeniedError struct {
	Message string
	Reason  string
}

func (e *AccessDeniedError) Error() string { return e.Message }

// ___________________________________________________________________________

// login error
type LoginError struct {
	Message string
	Reason  string
}

func (e *LoginError) Error() string { return e.Message }

func NewLoginError(message string) *LoginError { return &LoginError{Message: message} }

// ___________________________________________________________________________
