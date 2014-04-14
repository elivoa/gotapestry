package syd

import (
	"github.com/elivoa/got/config"
	"got/core"
	"got/utils"
)

// todo: think out a better way to register this.
var SYDModule = &core.Module{
	Name:            "syd",
	VarName:         "SYDModule",
	BasePath:        utils.CurrentBasePath(),
	PackagePath:     "syd",
	Description:     "SYD Selling System Main module.",
	IsStartupModule: true,
	Register: func() {
		c := config.Config

		// config static resources
		c.AddStaticResource("/pictures/", "/var/site/data/syd/pictures/")
		c.AddStaticResource("/static/", "../static/")
	},
}
