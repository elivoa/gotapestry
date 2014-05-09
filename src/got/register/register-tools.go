package register

import (
	"got/debug"
	"log"
	"reflect"
	"runtime"
)

// CreateURLs in it

/* ----------------------------------------------------------
   Tools
   ----------------------------------------------------------*/

var log_lookup bool = false

func logLookup(format string, params ...interface{}) {
	if log_lookup {
		log.Printf(format, params...)
	}
}

// ________________________________________________________________________________
// debug log
var debug_component_register = true

func debuglog(format string, params ...interface{}) {
	if debug_component_register {
		debug.Log(format, params...)
	}
}

// not used -----------------

// snippet.reflect: get function name
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
