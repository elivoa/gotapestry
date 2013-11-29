package person

import (
	"fmt"
	"got/core"
	"github.com/elivoa/gxl"
	"syd/model"
	"syd/service/personservice"
)

type EditAccountBallance struct {
	core.Page
	Id              *gxl.Int `path-param:"1"`
	Person          *model.Person
	AccountBallance float64
}

func (p *EditAccountBallance) Setup() {
	p.Person = personservice.GetPerson(p.Id.Int)
}

func (p *EditAccountBallance) OnSuccess() (string, string) {
	p.Setup() // init again
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", p.AccountBallance)
	if p.AccountBallance != p.Person.AccountBallance {
		p.Person.AccountBallance = p.AccountBallance
		fmt.Println(">>>>>>>>>>>>")
		if _, err := personservice.Update(p.Person); err != nil {
			fmt.Println(">>>>>>>>>>>>>>>>> ")
			panic(err)
		}
	}
	return "redirect", fmt.Sprintf("/person/EditAccountBallance/%v", p.Id.Int)
}
