package person

import (
	"github.com/elivoa/got/core"
	"github.com/elivoa/gxl"
	"syd/model"
	"syd/service"
)

type CustomerProfileCard struct {
	core.Component
	CustomerId *gxl.Int
	Customer   *model.Person
}

func (p *CustomerProfileCard) Setup() {
	if p.Customer == nil {
		if p.CustomerId == nil {
			panic("Customer or CustomerId should not both be null!")
		}
		// TODO get person
		var err error
		p.Customer, err = service.Person.GetPersonById(p.CustomerId.Int)
		if err != nil {
			panic(err)
		}
	}
}
