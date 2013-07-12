package personservice

import (
	"syd/dal/persondao"
	"syd/model"
)

func GetPerson(customerId int) *model.Customer {
	person, _ := persondao.Get(customerId)
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

func GetCustomer(customerId int) *model.Customer {
	person, _ := persondao.Get(customerId)
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

func List(personType string) ([]*model.Person, error) {
	return persondao.ListAll(personType)
}

func ListCustomer() ([]*model.Person, error) {
	return persondao.ListAll("customer")
}

func ListFactory() ([]*model.Person, error) {
	return persondao.ListAll("customer")
}

func Create(person *model.Person) (*model.Person, error) {
	return persondao.Create(person)
}
