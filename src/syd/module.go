package syd

import (
	"encoding/gob"
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/core/exception"
	"github.com/elivoa/got/coreservice/sessions"
	"github.com/elivoa/got/errorhandler"
	"github.com/elivoa/got/templates"
	"github.com/elivoa/got/utils"
	"github.com/elivoa/gxl"
	"net/http"
	"reflect"
	"strings"
	"syd/base"
	"syd/model"
)

var appversion = "3.5"

// todo: think out a better way to register this.
var SYDModule = &core.Module{
	Name:            "syd",       // Don't use this. It's only used to display.
	Version:         appversion,  // TODO: used to add to assets path to disable cache.
	VarName:         "SYDModule", // Variable name.
	BasePath:        utils.CurrentBasePath(),
	PackagePath:     "syd", // package name used anywhere to locate important things.
	Description:     "SYD Selling System Main module.",
	IsStartupModule: true,
	Register: func() {
		c := config.Config

		// config static resources
		c.AddStaticResource("/static/", "../static/") // static files.
		c.AddStaticResource("/ng/", "../ng/")         // angularjs modules.
		c.AddStaticResource("/pictures/", "/var/site/data/syd/pictures/")
		c.ResourcePath = "/var/site/data/syd/pictures/"

		c.Port = 8080 //13062 for server
		c.DBPort = 3306
		c.DBName = "syd"
		c.DBUser = "root"
		c.DBPassword = "eserver409$)("

		// Register new template functions this app usages.
		templates.RegisterFunc("HasAnyRole", HasAnyRole)
		templates.RegisterFunc("version", AppVersion)

		// Register error handlers.
		errorhandler.AddHandler("LoginError",
			reflect.TypeOf(base.LoginError{}),
			errorhandler.RedirectHandler("/account/login"),
		)

		errorhandler.AddHandler("TimeZoneNotFoundError",
			reflect.TypeOf(exception.TimeZoneNotFoundError{}),
			errorhandler.RedirectHandler("/account/login"),
		)

		// --------------------------------------------------------------------------------
		// Modify some default parmeters.
		config.LIST_PAGE_SIZE = 50
		config.ReloadTemplate = true // disable reload template?

		// Config gxl toolset.
		gxl.Locale = gxl.CN // set gxl toolset language to Chinese.
		gxl.Timezone = +0   // set gxl toolset timezone to 8 Beijing.

		// TODO: Register Coercers.

		// Register gob
		gob.Register(&model.UserToken{})

	},
}

// --------------------------------------------------------------------------------
// Additional functions.

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

func AppVersion() string {
	return appversion
}
