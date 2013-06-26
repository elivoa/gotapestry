package person

import (
	"fmt"
	"github.com/gorilla/mux"
	"got/core"
	"got/register"
	"gxl"
	"net/http"
	"strconv"
	"syd/dal"
	"syd/model"
)

func Register() {}
func init() {
	register.Page(Register,
		&PersonIndex{}, &PersonList{}, &PersonEdit{},
		&PersonDetail{},
	)
}

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
	Persons    *[]model.Person
	SubTitle   string
}

func (p *PersonList) Setup() interface{} {
	persons, err := dal.ListPerson(p.PersonType)
	if err != nil {
		return err
	}
	p.Persons = persons
	p.SubTitle = listTypeLabel[p.PersonType]
	return nil
	// return "template", "person-list"
}

// TODO
func (p *PersonList) Ondelete() {

}

/* ________________________________________________________________________________
   TODO
*/
func PersonDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	personType := vars["type"]

	// TODO: important: need some security validation
	dal.DeletePerson(id)

	// redirect to person list.
	http.Redirect(w, r, fmt.Sprintf("/person/list/%v", personType), http.StatusFound)
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
		p.Person = dal.GetPerson(p.Id.Int)
		p.SubTitle = "编辑"
	} else {
		p.Person = model.NewPerson()
		p.SubTitle = "新建"
	}
}

func (p *PersonEdit) OnSuccessFromPersonForm() (string, string) {
	if p.Id != nil {
		dal.UpdatePerson(p.Person)
	} else {
		dal.CreatePerson(p.Person)
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
	Orders *[]model.Order
}

func (p *PersonDetail) Setup() {
	if p.Id == nil {
		return
	}
	p.Person = dal.GetPerson(p.Id.Int)
	if p.Person != nil {
		p.Orders = dal.ListOrderByCustomer(p.Person.Id, "all")
	}
}
