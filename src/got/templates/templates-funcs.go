/*
   Time-stamp: <[templates-funcs.go] Elivoa @ Sunday, 2013-08-18 14:54:28>
*/
package templates

import (
	"gxl"
	"html/template"
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
	return gxl.FormatCurrency(d, 0)
}
