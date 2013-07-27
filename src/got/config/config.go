package config

import (
	"path"
	"runtime"
)

// ________________________________________________________________________________
// System configs

/*
 * TODO Auto-detect PagePackages
 * ...
 */
var Config = NewConfigure()

type Configure struct {
	Version      string `Framewrok Version`
	BasePath     string // /path/to/home
	SrcPath      string // /path/to/home/src
	StaticPath   string // /path/to/home/static
	ResourcePath string // /var/site/data/syd/

	BasePackages      []string `Packages that contains Pages and Components etc`
	PagePackages      []string `no use`
	ComponentPackages []string `...`

	// other system config
	TemplateFileExtension string
}

func NewConfigure() Configure {
	// 1. get base path
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic("Can't get current path!")
	}
	basePath := path.Join(path.Dir(file), "../../..")

	return Configure{
		Version:      "Alpha 3",
		BasePackages: []string{"happystroking"},
		BasePath:     basePath,
		ResourcePath: "/tmp/",

		SrcPath:    path.Join(basePath, "src"),
		StaticPath: path.Join(basePath, "static"),

		TemplateFileExtension: ".html",
	}
}
