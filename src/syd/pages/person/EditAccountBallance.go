package person

import (
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/gxl"
	"syd/dal/accountdao"
	"syd/model"
	"syd/service"
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
	var err error
	p.Person, err = service.Person.GetPersonById(p.Id.Int)
	if err != nil {
		panic(err)
	}
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
