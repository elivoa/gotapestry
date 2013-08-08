/*
   Time-stamp: <[templates-component.go] Elivoa @ Thursday, 2013-08-08 19:10:28>
*/

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

// register components as template function call.
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
