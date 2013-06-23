package got

import (
	"fmt"
	"got/core"
	"got/register"
)

func init() {
	register.Page(Register, &Errors{})
}

func Register() {}

type Errors struct {
	core.Page
}

func (p *Errors) SetupRender() {
	fmt.Println("\n\nPage Error page")
}
