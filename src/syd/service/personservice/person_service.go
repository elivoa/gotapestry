package personservice

import (
	"syd/dal"
	"syd/model"
)

func GetCustomer(customerId int) *model.Customer {
	person := dal.GetPerson(customerId)
	if nil == person {
		return nil
	}
	customer := model.Customer{
		Person:      *person,
		Accumulated: 998,
	}
	// TODO get Accumulated
	return &customer
}
