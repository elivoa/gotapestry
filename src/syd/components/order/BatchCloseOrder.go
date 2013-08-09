package order

import (
	"got/core"
)

type BatchCloseOrder struct {
	core.Component
	CustomerId int // TODO use this.
	// Customer   *model.Person // now use this.
	Class            string //  link style
	JSInit           bool   // false if you want to manually init js
	QuickClearButton bool   // 已结清按钮
}

func (c *BatchCloseOrder) New() *BatchCloseOrder {
	return &BatchCloseOrder{
		JSInit:           true,
		QuickClearButton: true,
	}
}

func (c *BatchCloseOrder) Setup() {
	// fmt.Println("------------------------", c.Customer.AccountBallance)
}
