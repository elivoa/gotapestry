package got

import (
	"fmt"
	"got/core"
)

type Errors struct {
	A *string
	core.Page
	C []int
}

func (p *Errors) SetupRender() {
	fmt.Println("\n\nPage Error page")
}
