/*
   Time-stamp: <[templates-funcs.go] Elivoa @ Monday, 2013-07-29 12:12:56>
*/
package templates

import (
	"fmt"
	"html/template"
	"time"
)

// TODO open this to developer to register global functions.
func registerBuiltinFuncs() {
	// init functions
	Templates.Funcs(template.FuncMap{
		// deprecated
		"eq": equas,

		"beautytime":     BeautyTime,
		"beautycurrency": BeautyCurrency,

		// new
		"formattime":     FormatTime,
		"formatcurrency": BeautyCurrency,
		"prettycurrency": BeautyCurrency,
		"prettytime":     BeautyTime,
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

func BeautyCurrency(d float64) string {
	return fmt.Sprintf("%.2f", d)
}
