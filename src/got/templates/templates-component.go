package templates

import (
	// "encoding/json"
	"fmt"
	"html/template"
	"log"
	"strings"
)

/*_______________________________________________________________________________
  Register components
*/
func registerComponentFuncs() {
	// init functions
	Templates.Funcs(template.FuncMap{})
}

func RegisterComponent(name string, f interface{}) {
	funcName := fmt.Sprintf("t_%v", strings.Replace(name, "/", "_", -1))
	debuglog("-108- [RegisterComponent] %v", funcName)
	Templates.Funcs(template.FuncMap{funcName: f})
}

// --------------------------------------------------------------------------------
// log
//
var debugLog = true

func debuglog(format string, params ...interface{}) {
	if debugLog {
		log.Printf(format, params...)
	}
}
