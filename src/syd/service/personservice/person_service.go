package personservice

import (
	"sort"
	"syd/dal/persondao"
	"syd/model"
	"syd/service/suggest"
)

// is this correct?
func GetProducer(id int) *model.Producer {
	person, _ := persondao.Get(persondao.EntityManager().PK, id)
	if nil == person {
		return nil
	}
	producer := model.Producer{
		Person: *person,
	}
	return &producer
}

// is this correct?
func GetCustomer(customerId int) *model.Customer {
	person, _ := persondao.Get(persondao.EntityManager().PK, customerId)
	if nil == person {
		return nil
	}
	customer := model.Customer{
		Person: *person,
	}
	return &customer
}

// --------------------------------------------------------------------------------

// person list, sortable
type PersonListSortbyAccountBallance []*model.Person

func (p PersonListSortbyAccountBallance) Len() int { return len(p) }
func (p PersonListSortbyAccountBallance) Less(i, j int) bool {
	return p[i].AccountBallance < p[j].AccountBallance
}
func (p PersonListSortbyAccountBallance) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func SortByAccumulated(list []*model.Person) {
	newl := PersonListSortbyAccountBallance{}
	newl = list[:]
	sort.Sort(newl)
}

// // person list, sortable
// type PersonListSortbyAccountBallance struct {
// 	list []*model.Person
// }

// func (p *PersonListSortbyAccountBallance) Len() int {
// 	return len(p.list)
// }

// func (p *PersonListSortbyAccountBallance) Less(i, j int) bool {
// 	return p.list[i].AccountBallance < p.list[j].AccountBallance
// }

// func (p *PersonListSortbyAccountBallance) Swap(i, j int) {
// 	t := p.list[j]
// 	p.list[j] = p.list[i]
// 	p.list[i] = t
// }

// func SortByAccumulated(list []*model.Person) *PersonListSortbyAccountBallance {
// 	newl := &PersonListSortbyAccountBallance{
// 		list: list,
// 	}
// 	sort.Sort(newl)
// 	return newl
// }

// --------------------------------------------------------------------------------

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

// Note: Don't forget create an AccountChangeLog when update person's account-ballance.
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
