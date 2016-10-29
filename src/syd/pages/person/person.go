package person

import (
	"fmt"
	"github.com/elivoa/got/core"
	"syd/model"
	"syd/service/personservice"
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
