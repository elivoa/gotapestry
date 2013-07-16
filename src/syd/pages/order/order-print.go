package order

import (
	"fmt"
	"got/core"
	"got/register"
)

func init() {
	register.Page(Register, &OrderPrint{})
}

// ________________________________________________________________________________
// OrderPrint

type OrderPrint struct {
	core.Page
	TrackNumber string `path-param:"1"`
}

func (p *OrderPrint) Activate() {
	fmt.Println("Order Print is here")
}
