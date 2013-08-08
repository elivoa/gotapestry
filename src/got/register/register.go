package register

import (
	"got/core"
)

var (
	Pages      = ProtonSegment{Name: "/"}
	Components = ProtonSegment{Name: "/"}
	Mixins     = ProtonSegment{Name: "/"}
)

// ----  Register Pages  ----------------------------------------------------------------------------

// TODO delete this.

// Register a Page.
// Params
//   - f used to locate a page(reflect has no use).
//   - pages are Pages located in that folder.
func Page(f func(), pages ...core.Pager) int {
	// fmt.Println("some one register a page", pages)
	// for _, p := range pages {
	// 	url := MakeUrl(f, p)
	// 	Pages.Add(url, p)
	// }
	// return len(pages)
	return 0
}

// ----  Register Components  ----------------------------------------------------------------------------
