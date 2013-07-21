package route

import (
	"strings"
)

func RedirectDispatch(targets ...string) (string, string) {
	for _, target := range targets {
		if strings.Trim(target, " ") != "" {
			return "redirect", target
		}
	}
	panic("Can't Dispatch any of these redirects.")
	return "error", "Can't Dispatch any of these redirects."
}
