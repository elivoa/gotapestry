package components

import (
	"got/register"
)

func Register() {} // Used to locate package, any better approach?

func init() {
	register.Component(Register, &SuggestControl{}, &OrderDetailsEditor{})
}
