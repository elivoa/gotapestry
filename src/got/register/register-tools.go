package register

import (
	"got/debug"
	"log"
	"reflect"
	"runtime"
)

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

// // join two part of url.
// // TODO performance here is bad. design is bad.
// func MakeUrl(f func(), p interface{}) string {
// 	fName := GetFunctionName(f)
// 	prefix := fName[:strings.LastIndex(fName, ".")]
// 	suffix := reflect.TypeOf(p).String()
// 	suffix = suffix[strings.LastIndex(suffix, ".")+1:]
// 	url := fmt.Sprintf("%v/%v", prefix, suffix)
// 	return url
// }

// not used -----------------

// snippet.reflect: get function name
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
