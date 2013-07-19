package templates

import (
	"fmt"
	"html/template"
	"time"
)

func registerHelperFuncs() {
	// init functions
	Templates.Funcs(template.FuncMap{
		"formattime":     FormatTime,
		"beautytime":     BeautyTime,
		"formatcurrency": BeautyCurrency,
		"beautycurrency": BeautyCurrency,
	})
}

/*_______________________________________________________________________________
  Tempalte Functions
*/

// {{showtime .CreateTime "2006-01-02 15:04:05"}}
func FormatTime(t time.Time, format string) string {
	return t.Format(format)
}

func BeautyTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func BeautyCurrency(d float64) string {
	return fmt.Sprintf("%.2f", d)
}
