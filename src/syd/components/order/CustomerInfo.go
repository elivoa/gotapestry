package order

import (
	"fmt"
	"got/core"
	"syd/model"
)

type CustomerInfo struct {
	core.Component
	Customer    *model.Customer
	Accumulated float64
}

func (c *CustomerInfo) Setup() {
	// if c.Accumulated is not set, use customer's
	if !c.Injected("Accumulated") {
		c.Accumulated = -c.Customer.AccountBallance
		fmt.Println("================================= set accumulated to ", c.Customer.AccountBallance)
	}
}
