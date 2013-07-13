package personservice

import (
	"syd/dal/persondao"
	"syd/model"
)

// is this correct?
func GetProducer(id int) *model.Producer {
	person, _ := persondao.Get(id)
	if nil == person {
		return nil
	}
	producer := model.Producer{
		Person: *person,
		// Accumulated: 998,
	}
	return &producer
}

// is this correct?
func GetPerson(id int) *model.Person {
	person, _ := persondao.Get(id)
	if nil == person {
		return nil
	}
	return person
}

// is this correct?
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

// only return error?
func Create(person *model.Person) error {
	return persondao.Create(person)
}

func Update(person *model.Person) (int64, error) {
	return persondao.Update(person)
}

func Delete(id int) (int64, error) {
	return persondao.Delete(id)
}
