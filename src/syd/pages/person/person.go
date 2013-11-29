package person

import (
	"fmt"
	"got/core"
	"github.com/elivoa/gxl"
	"syd/dal/orderdao"
	"syd/dal/persondao"
	"syd/model"
	"syd/service/orderservice"
	"syd/service/personservice"
	"time"
)

var (
	listTypeLabel = map[string]string{"customer": "客户", "factory": "厂商"}
)

/*_______________________________________________________________________________
  Person List Page
*/
type PersonIndex struct{ core.Page }

func (p *PersonIndex) Setup() (string, string) {
	return "redirect", "/person/list/customer"
}

/*_______________________________________________________________________________
  Person List Page
*/
type PersonList struct {
	core.Page

	PersonType string `path-param:"1" param:"type"`
	Persons    []*model.Person
	SubTitle   string
}

func (p *PersonList) Setup() interface{} {
	persons, err := personservice.List(p.PersonType)
	if err != nil {
		return err
	}
	p.Persons = persons
	p.SubTitle = listTypeLabel[p.PersonType]
	return nil
	// return "template", "person-list"
}

func (p *PersonList) Ondelete(personId int, personType string) (string, string) {
	personservice.Delete(personId)
	// TODO make this default redirect.
	return "redirect", fmt.Sprintf("/person/list/%v", personType)
}

/*_______________________________________________________________________________
  Person Edit Page
*/
type PersonEdit struct {
	core.Page

	Id     *gxl.Int `path-param:"1" required:"true" param:"id"`
	Person *model.Person

	Title    string
	SubTitle string

	TypeData interface{} // for type select
}

func (p *PersonEdit) Activate() {
	// here is some lightweight init.
	p.TypeData = &listTypeLabel
}

func (p *PersonEdit) Setup() {
	p.Title = "create/edit Person"

	if p.Id != nil {
		person, err := persondao.Get(p.Id.Int)
		if err != nil {
			// TODO how to handle error on page object?
			panic(err.Error())
		}
		p.Person = person
		p.SubTitle = "编辑"
	} else {
		p.Person = model.NewPerson()
		p.SubTitle = "新建"
	}
}

func (p *PersonEdit) OnSubmit() {
	if p.Id != nil {
		p.Person = personservice.GetPerson(p.Id.Int)
	} else {
		// No Need to edit.
	}
}

func (p *PersonEdit) OnSuccess() (string, string) {
	if p.Id != nil {
		personservice.Update(p.Person)
	} else {
		personservice.Create(p.Person)
	}
	return "redirect", fmt.Sprintf("/person/list/%v", p.Person.Type)
}

/* ________________________________________________________________________________
   PersonEdit
*/
type PersonDetail struct {
	core.Page

	Id *gxl.Int `path-param:"1"`

	Person *model.Person
	Orders []*model.Order
	// TheBigOrder    *model.Order
	// LeavingMessage string
	TodayOrders []*model.Order
}

func (p *PersonDetail) LeavingMessage(order *model.Order) string {
	return orderservice.LeavingMessage(order)
}

func (p *PersonDetail) Setup() {
	if p.Id == nil {
		return
	}
	// performance issue: here we load all orders, this has an performance issue.
	p.Person = personservice.GetPerson(p.Id.Int)
	if p.Person != nil {
		orders, err := orderdao.ListOrderByCustomer(p.Person.Id, "all")
		if err != nil {
			panic(err.Error())
		}
		p.Orders = orders
	}

	// get today orders.
	date := time.Now()
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 0, 1)
	orders, err := orderdao.ListOrderByCustomer_Time(p.Person.Id, start, end)
	if err != nil {
		panic(err.Error())
	}
	orderservice.LoadDetails(orders)
	p.TodayOrders = orders

	// p.TheBigOrder, p.LeavingMessage = orderservice.GenerateLeavingMessage(p.Person.Id, time.Now())
	if true {
		return
	}
}
