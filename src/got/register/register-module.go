package register

import (
	"fmt"
	"path"
	"sync"
)

// Use Module instead.
// TODO:
//   - DONE| Cache module
//   - Remove APP
//
var Modules = &ModuleCache{m: map[string]*Module{}}

type ModuleCache struct {
	l sync.RWMutex
	m map[string]*Module
}

func RegisterModule(modules ...*Module) {
	for _, m := range modules {
		Modules.Add(m)
	}
}

// e.g: for /Users/xxx/src/module/pages/xxx
type Module struct {
	Name        string // module name
	BasePath    string // full file path of this module. e.g.: /Users/xxx/src/
	PackagePath string // package path. e.g.: /module
	Description string

	Register func() // manually register page and components.
}

// Path returns /User/xxx/src/module
func (m *Module) Path() string {
	return path.Join(m.BasePath, m.PackagePath)
}

func (mc *ModuleCache) Add(module *Module) {
	mc.l.Lock()
	mc.m[module.Name] = module
	mc.l.Unlock()
}

func (mc *ModuleCache) Get(name string) *Module {
	mc.l.RLock()
	module := mc.m[name]
	mc.l.RUnlock()
	return module
}

func (mc *ModuleCache) Map() map[string]*Module {
	return mc.m
}

// ----  Printing  -----------------------------------------------------------------------------------

func (m *Module) String() string {
	return m.Path()
}

func (mc *ModuleCache) PrintALL() {
	fmt.Println("---- Modules ---------------------")
	mc.l.RLock()
	for _, module := range mc.m {
		fmt.Printf("  %v\n", module.String())
	}
	mc.l.RUnlock()
}

// // ---------------------------------------------------------------------------------------------------
// // TODO delete all tings below...

// ----  Advanced Functions  -------------------------------------------------------------------------

// GetPaths returns all module's file path, used to parse files.
// func (mc *ModuleCache) GetPaths() []string {
// 	mc.l.RLock()
// 	defer mc.l.RUnlock()

// 	paths := make([]string, len(mc.m))
// 	i := 0
// 	for _, module := range mc.m {
// 		paths[i] = module.Path
// 		i += 1
// 	}
// 	return paths
// }
