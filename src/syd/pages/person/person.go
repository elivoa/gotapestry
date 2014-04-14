package person

import (
	"fmt"
	"github.com/elivoa/gxl"
	"got/core"
	"syd/dal/accountdao"
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

	Current   int `path-param:"2"` // pager: the current item. in pager.
	PageItems int `path-param:"3"` // pager: page size.
	Total     int

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

	// fix pagers
	if p.PageItems <= 0 {
		p.PageItems = 50 // TODO default pager number. Config this.
	}

	// performance issue: here we load all orders, this has an performance issue.
	p.Person = personservice.GetPerson(p.Id.Int)
	if p.Person != nil {
		var err error
		p.Total, err = orderdao.CountOrderByCustomer("all", p.Person.Id)
		if err != nil {
			panic(err.Error())
		}
		// TODO finish the common conditional query.
		// query with pager
		p.Orders, err = orderdao.ListOrderByCustomer(p.Person.Id, "all", p.Current, p.PageItems)
		if err != nil {
			panic(err.Error())
		}
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

func (p *PersonDetail) UrlTemplate() string {
	return fmt.Sprintf("/person/detail/%d/{{Start}}/{{PageItems}}", p.Person.Id)
}

func (p *PersonDetail) ShouldShowLeavingMessage(o *model.Order) bool {
	switch model.OrderType(o.Type) {
	// case model.Wholesale, model.ShippingInstead:
	case model.Wholesale, model.SubOrder:
		return true
	}
	return false
}

func (p *PersonDetail) ChangeLogs() []*model.AccountChangeLog {
	changes, err := accountdao.ListAccountChangeLogsByCustomerId(p.Person.Id)
	if err != nil {
		panic(err.Error())
	}
	return changes
}

func (p *PersonDetail) DisplayType(t int) string {
	switch t {
	case 1:
		return "手工修改"
	case 2:
		return "订单发货"
	case 3:
		return "批量结款"
	case 4:
		return "取消已发货订单，减去累计欠款"
	default:
		return "不可知类型"
	}
}
