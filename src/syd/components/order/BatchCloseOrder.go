package order

import (
	"got/core"
	"got/register"
	"syd/model"
	"fmt"
)

func init() {
	register.Component(Register, &BatchCloseOrder{})
}

type BatchCloseOrder struct {
	core.Component
	CustomerId int           // TODO use this.
	Customer   *model.Person // now use this.
}

func (c *BatchCloseOrder) Setup() {
	fmt.Println("------------------------", c.Customer.AccountBallance)
}
