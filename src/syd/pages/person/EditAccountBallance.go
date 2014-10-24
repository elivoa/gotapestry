package person

import (
	"fmt"
	"github.com/elivoa/gxl"
	"github.com/elivoa/got/core"
	"syd/dal/accountdao"
	"syd/model"
	"syd/service/personservice"
)

type EditAccountBallance struct {
	core.Page
	Id              *gxl.Int `path-param:"1"`
	Person          *model.Person
	AccountBallance float64
	Reason          string
}

func (p *EditAccountBallance) Setup() {
	p.Person = personservice.GetPerson(p.Id.Int)
}

// on form submit
func (p *EditAccountBallance) OnSuccess() (string, string) {
	// init person again. get person.
	p.Setup()

	if p.AccountBallance != p.Person.AccountBallance {
		// create
		accountdao.CreateAccountChangeLog(&model.AccountChangeLog{
			CustomerId:     p.Person.Id,
			Delta:          p.AccountBallance - p.Person.AccountBallance,
			Account:        p.AccountBallance,
			Type:           1, // manually modification;
			RelatedOrderTN: 0,
			Reason:         p.Reason,
		})
		// update account ballance
		p.Person.AccountBallance = p.AccountBallance
		if _, err := personservice.Update(p.Person); err != nil {
			panic(err)
		}
	}
	return "redirect", fmt.Sprintf("/person/EditAccountBallance/%v", p.Id.Int)
}
