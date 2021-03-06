package order

import (
	"github.com/elivoa/got/core"
	"strconv"
	"syd/model"
	"syd/service"
)

type BatchCloseOrder struct {
	core.Component
	CustomerId       int           // TODO use this.
	Customer         *model.Person // now use this.
	Class            string        //  link style
	JSInit           bool          // false if you want to manually init js
	QuickClearButton bool          // 已结清按钮
}

func (c *BatchCloseOrder) New() *BatchCloseOrder {
	return &BatchCloseOrder{
		JSInit:           true,
		QuickClearButton: true,
	}
}

func (c *BatchCloseOrder) Setup() {
	var err error
	c.Customer, err = service.Person.GetPersonById(c.CustomerId)
	if err != nil {
		panic(err)
	}
	if c.Customer == nil {
		panic("Customer not found! " + strconv.Itoa(c.CustomerId))
	}
}
