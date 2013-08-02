package components

import (
	"got/route"
)

func Register() {} // Used to locate package, any better approach?

func init() {
	route.Component(Register, &SuggestControl{}, &OrderDetailsEditor{})
}
