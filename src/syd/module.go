package syd

import (
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/utils"
	"got/core"
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
