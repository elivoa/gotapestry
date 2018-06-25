package person

import (
	"fmt"
	"strconv"
	"syd/model"
	"syd/service"
	"syd/service/personservice"

	"github.com/elivoa/got/core"
	"github.com/elivoa/gxl"
)

type PersonEdit struct {
	core.Page

	Id     *gxl.Int `path-param:"1" required:"true" param:"id"`
	Person *model.Person

	Title    string
	SubTitle string

	// bool options
	IsSendNewProduct bool // 是否发样衣
	IsPrintHidePrice bool // 是否默认不打印价格

	TypeData  interface{} // for type select
	LevelData interface{} // for type select
	HideData  interface{} // for type select
}

func (p *PersonEdit) Activate() {
	// here is some lightweight init.
	p.TypeData = &listTypeLabel
	p.LevelData = &levelTypeLabel
	p.HideData = &hideTypeLabel
}

func (p *PersonEdit) Setup() {
	p.Title = "create/edit Person"

	if p.Id != nil {
		person, err := service.Person.GetPersonById(p.Id.Int)
		if err != nil {
			// TODO how to handle error on page object?
			panic(err.Error())
		}
		p.Person = person
		p.SubTitle = "编辑"

		// Extend, load bool options
		{
			value, err := service.Const.Get2ndIntValue("SendNewProduct", strconv.Itoa(person.Id))
			if err != nil {
				panic(err)
			}
			p.IsSendNewProduct = value >= 0

			value2, err := service.Const.Get2ndIntValue("PrintHidePrice", strconv.Itoa(person.Id))
			if err != nil {
				panic(err)
			}
			p.IsPrintHidePrice = value2 >= 0
		}

	} else {
		p.Person = model.NewPerson()
		p.SubTitle = "新建"
	}
}

func (p *PersonEdit) OnPrepareForSubmit() {
	if p.Id != nil {
		var err error
		p.Person, err = service.Person.GetPersonById(p.Id.Int)
		if err != nil {
			panic(err)
		}
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
