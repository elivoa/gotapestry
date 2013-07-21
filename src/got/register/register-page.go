package register

import (
	"got/core"
)

/* ________________________________________________________________________________
   PageRegister
*/
var Pages = ProtonSegment{Name: "/"}

/* ________________________________________________________________________________
   Register a Page,
   params:
     - f used to locate a page(reflect has no use).
     - pages are Pages located in that folder.
*/
func Page(f func(), pages ...core.Pager) int {
	for _, p := range pages {
		url := makeUrl(f, p)
		Pages.Add(url, p, "page")
	}
	return len(pages)
}
