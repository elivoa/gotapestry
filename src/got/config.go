package got

import ()

// ________________________________________________________________________________
// System configs

/*
 * TODO Auto-detect PagePackages
 * ...
 */
var Config = NewConfigure()

type Configure struct {
	Version           string   `Framewrok Version`
	BasePackages      []string `Packages taht contains Pages and Components etc`
	PagePackages      []string `no use`
	ComponentPackages []string `...`
}

func NewConfigure() Configure {
	return Configure{
		Version:      "Alpha 2",
		BasePackages: []string{"happystroking"},
	}
}

// ________________________________________________________________________________
// Got start configs

type GotConfig struct {
	StaticResources []string // e.g.: ["/static/", "../"]
}
