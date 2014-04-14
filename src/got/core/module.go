package core

import (
	"path"
)

// e.g: for /Users/xxx/src/module/pages/xxx
/*
  Notes:
    {{.PackagePath}}.{{.VarName}} // calling methods

*/
type Module struct {
	Name            string // seems no use.
	VarName         string // module name, should be the same with struct name.
	BasePath        string // full file path of this module. e.g.: /Users/xxx/src/
	PackagePath     string // package path. e.g.: /module
	Description     string
	IsStartupModule bool

	// Register config the module in more details.
	// This method only called by generated code.
	Register func() // manually register page and components.
}

// Path returns /User/xxx/src/module
func (m *Module) Path() string {
	return path.Join(m.BasePath, m.PackagePath)
}

func (m *Module) String() string {
	return m.Path()
}

/* Example */
var ___Example_Module___ Module = Module{
	Name:            "syd",
	VarName:         "SYDModule",
	BasePath:        "/Users/bogao/develop/gitme/gotapestry/src",
	PackagePath:     "syd",
	Description:     "SYD Selling System Main module.",
	IsStartupModule: true,
}
