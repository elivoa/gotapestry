package syd

import (
	"got/config"
	"got/register"
	"got/utils"
)

// todo: think out a better way to register this.
var SYDModule = &register.Module{
	Name:        "syd",
	BasePath:    utils.CurrentBasePath(),
	PackagePath: "syd",
	Description: "SYD Selling System Main module.",
	Register: func() {
		c := config.Config

		// config static resources
		c.AddStaticResource("/pictures/", "/var/site/data/syd/pictures/")
		c.AddStaticResource("/static/", "../static/")
	},
}
