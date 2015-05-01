package admin

import (
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route"
	"syd/dal/constdao"
	"syd/model"
	"syd/service"
)

type ConstList struct {
	core.Component
	Name     string
	Consts   []*model.Const
	HideName bool

	Referer  string
	TimeZone *model.TimeZoneInfo
}

func (p *ConstList) Setup() interface{} {
	// service.User.RequireRole(p.W, p.R, carfilm.RoleSet_Management...)
	// p.TimeZone = service.TimeZone.UserTimeZoneSafe(p.R)

	// load it.
	consts, err := constdao.GetList(p.Name)
	if err != nil {
		panic(err.Error())
	}
	p.Consts = consts

	// fmt.Println("\n\n\n\n$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	// fmt.Println("things i want to see is:  ", p.Referer)

	return true
}

func (p *ConstList) Ondelete(id int64) interface{} {
	fmt.Println("\n\n\n\n$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	fmt.Println("things i want to see is:  ", p.Referer)
	if _, err := service.Const.DeleteById(id); err != nil {
		panic(err)
	}
	// return exit.RedirectFirstValid(p.Referer, "/admin/preference/")
	refererurl := route.GetRefererFromURL(p.Request())
	return route.RedirectDispatch(refererurl, "/admin/preference/")
}
