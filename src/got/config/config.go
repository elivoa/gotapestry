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
	Version      string // `Framewrok Version`
	BasePath     string // /path/to/home
	SrcPath      string // /path/to/home/src
	StaticPath   string // /path/to/home/static
	ResourcePath string // /var/site/data/syd/

	// module path need to import. this is not yours.
	ModulePath      []*ModulePath // [[ImportPath, StructName], ...]
	StaticResources [][]string    // e.g.: [["/static/", "../"], ...]

	BasePackages      []string `Packages that contains Pages and Components etc`
	PagePackages      []string `no use`
	ComponentPackages []string `...`

	// other system config
	TemplateFileExtension string

	// server
	Port int // start port
}

func NewConfigure() *Configure {
	// 1. get base path
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic("Can't get current path!")
	}
	basePath := path.Join(path.Dir(file), "../../..")

	return &Configure{
		Version:      "0.2.0",
		BasePackages: []string{"happystroking"},
		BasePath:     basePath,
		ResourcePath: "/tmp/",

		SrcPath:    path.Join(basePath, "src"),
		StaticPath: path.Join(basePath, "static"),

		ModulePath:            []*ModulePath{},
		StaticResources:       [][]string{},
		TemplateFileExtension: ".html",

		// server
		Port: 8080,
	}
}

func (c *Configure) RegisterModulePath(importPath string, name string) {
	Config.ModulePath = append(Config.ModulePath, &ModulePath{PackagePath: importPath, Name: name})
}

func (c *Configure) AddStaticResource(url string, path string) {
	// TODO warn | after log
	Config.StaticResources = append(Config.StaticResources, []string{url, path})
}

type ModulePath struct {
	PackagePath string
	Name        string
}
