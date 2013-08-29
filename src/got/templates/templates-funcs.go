/*
   Time-stamp: <[templates-funcs.go] Elivoa @ Monday, 2013-08-26 17:06:57>
*/
package templates

import (
	"gxl"
	"html/template"
	"math"
	"time"
)

// TODO open this to developer to register global functions.
func registerBuiltinFuncs() {
	// init functions
	Templates.Funcs(template.FuncMap{
		// deprecated
		"eq": equas,

		// new
		"formattime":     FormatTime,
		"prettytime":     BeautyTime,
		"prettyday":      gxl.PrettyDay,
		"prettycurrency": PrettyCurrency,
		
		"now":            func() time.Time { return time.Now() },
	})
}

/*_______________________________________________________________________________
  Tempalte Functions
*/

func equas(o1 interface{}, o2 interface{}) bool {
	return o1 == o2
}

// {{showtime .CreateTime "2006-01-02 15:04:05"}}
func FormatTime(format string, t time.Time) string {
	return t.Format(format)
}

func BeautyTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func PrettyCurrency(d float64) string {
	if math.Mod(d, 1) > 0 {
		return gxl.FormatCurrency(d, 2)
	} else {
		return gxl.FormatCurrency(d, 0)
	}
}
