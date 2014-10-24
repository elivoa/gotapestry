package person

import (
	"github.com/elivoa/got/core"
	"syd/model"
	"syd/service/personservice"
)

type DebtCustomerList struct {
	core.Component
	Customers      []*model.Person
	SumAccumulated float64
}

func (c *DebtCustomerList) SetupRender() {
	customers, err := personservice.ListCustomer()
	if err != nil {
		panic(err.Error())
	}
	personservice.SortByAccumulated(customers)
	for i := 0; i < len(customers); i++ {
		ab := customers[i].AccountBallance
		if ab == 0 {
			customers[i] = nil
		} else if ab < 0 {
			c.SumAccumulated += -ab
		}
	}
	c.Customers = customers
}
