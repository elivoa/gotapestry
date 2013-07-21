package personservice

import (
	"syd/dal/persondao"
	"syd/model"
	"syd/service/suggest"
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
	if err := persondao.Create(person); err != nil {
		return err
	}
	// update suggest
	if person.Type == "customer" {
		suggest.Add(suggest.Customer, person.Name, person.Id)
	} else if person.Type == "factory" {
		suggest.Add(suggest.Factory, person.Name, person.Id)
	}
	return nil
}

func Update(person *model.Person) (affacted int64, err error) {
	if affacted, err = persondao.Update(person); err != nil {
		return
	}
	// update person
	if person.Type == "customer" {
		suggest.Update(suggest.Customer, person.Name, person.Id)
	} else if person.Type == "factory" {
		suggest.Update(suggest.Factory, person.Name, person.Id)
	}
	return
}

func Delete(id int) (affacted int64, err error) {
	if affacted, err = persondao.Delete(id); err != nil {
		return
	}
	// update person ( factory and customer's id is nocollapse)
	suggest.Delete(suggest.Customer, id)
	suggest.Delete(suggest.Factory, id)
	return
}
