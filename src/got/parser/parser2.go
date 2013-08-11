/*
  Parse all page/components source files. Cache it's content.

  TODO:
    - Rename function names.
    - Cleanup this file.
    - Make a lightwight parse, only parse modules with core embed.
*/
package parser

import (
	"got/core"
	"got/utils"
	"strings"
	"sync"
)

/*
   SourceInfo -> --
   TypeInfo -> StructInfo

*/

// this file is new version of parser.go. with all things myself.
// SourceInfo is the top-level struct containing all extracted information
// about the app source code, used to generate main.go.
// TODO rename
type SourceInfo struct {
	// StructSpecs lists type info for all structs found under the code paths.
	// They may be queried to determine which ones (transitively) embed certain types.
	Structs []*StructInfo

	// map: importpath.name -> StructInfo; e.g.:
	StructMap map[string]*StructInfo

	//# TODO parse validation keys.
	// ValidationKeys provides a two-level lookup.  The keys are:
	// 1. The fully-qualified function name,
	//    e.g. "github.com/robfig/revel/samples/chat/app/controllers.(*Application).Action"
	// 2. Within that func's file, the line number of the (overall) expression statement.
	//    e.g. the line returned from runtime.Caller()
	// The result of the lookup the name of variable being validated.
	ValidationKeys map[string]map[int]string

	//# TODO refactor this.
	// A list of import paths.
	// R notices files with an init() function and imports that package.
	//// what's this? why this?
	InitImportPaths []string

	//#TODO kill this?
	// controllerSpecs lists type info for all structs found under
	// app/controllers/... that embed (directly or indirectly) revel.Controller
	// controllerSpecs []*TypeInfo

	// testSuites list the types that constitute the set of application tests.
	// testSuites []*TypeInfo

	l sync.RWMutex
}

// TypeInfo summarizes information about a struct type in the app source code.
// Cache everything of a struct.
type StructInfo struct {
	StructName    string    // e.g. "PageStructName"
	ImportPath    string    // e.g. "github.com/elivoa/app/pages/admin"
	PackageName   string    // e.g. "admin"
	ModulePackage string    // e.g. "got/builtin, syd"
	FilePath      string    // full file path.
	ProtonKind    core.Kind // + Page | Component | Mixins
	MethodSpecs   []*MethodSpec

	// Used internally to identify controllers that indirectly embed *revel.Controller.
	embeddedTypes []*embeddedTypeName
}

// IsProton returns true if is a Page, Component or Mixins. Must be exported struct.
func (t *StructInfo) IsProton() bool {
	switch t.ProtonKind {
	case core.PAGE, core.COMPONENT, core.MIXIN:
		return utils.IsCapitalized(t.StructName) // must be exported struct
	}
	return false
}

// ProtonPath returns `proton path`, e.g.: order/list
func (t *StructInfo) ProtonPath() string {
	pp := t.ImportPath[len(t.ModulePackage)+1:]
	switch {
	case strings.HasPrefix(pp, "pages"):
		return pp[5:]
	case strings.HasPrefix(pp, "components"):
		return pp[10:]
	case strings.HasPrefix(pp, "mixins"):
		return pp[6:]
	}
	return pp
}

func (t *StructInfo) PrintEmbedTypes() {
	if t.embeddedTypes != nil {
		for _, e := range t.embeddedTypes {
			e.String()
		}
	}
}
