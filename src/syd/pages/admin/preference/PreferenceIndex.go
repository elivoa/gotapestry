package preference

import (
	"fmt"
	"github.com/elivoa/got/core"
	"strings"
	"syd/dal/constdao"
	"syd/model"
	"syd/service"
)

type PreferenceIndex struct {
	core.Page
	Title string

	// parameters
	Consts  []*model.Const
	Tab     string `path-param:"1"` // tab: today, yestoday, all
	Edit    int64  `query:"edit"`   // const id
	Referer string `query:"referer"`

	// for create
	Name  string
	Const *model.Const // for create.

	// Current       int    `path-param:"2"` // pager: the current item. in pager.
	// PageItems     int    `path-param:"3"` // pager: page size.

	// properties
	// Total int // pager: total items available
}

func (p *PreferenceIndex) Activate() {
	// service.User.RequireRole(p.W, p.R, carfilm.RoleSet_Management...)
	if p.Tab == "" {
		p.Tab = "today" // default go in toprint
	}
}

func (p *PreferenceIndex) SetupRender() {
	p.Title = "Admin | Preference"

	if p.Edit > 0 {
		constModel, err := constdao.GetById(p.Edit)
		if err != nil {
			fmt.Printf("Err is %v\n", err) // panic(err)
		}
		p.Const = constModel
	} else {
		p.Const = new(model.Const)
	}
}

func (p *PreferenceIndex) TabStyle(tab string) string {
	if strings.ToLower(p.Tab) == strings.ToLower(tab) {
		return "cur"
	}
	return ""
}

func (p *PreferenceIndex) OnSuccessFromCreateConstForm() /* *exit.Exit  */ {
	// service.User.RequireRole(p.W, p.R, carfilm.RoleSet_Management...)

	if p.Const.Id > 0 {
		// update by id
		err := service.Const.Update(p.Const.Name, p.Const.Key, p.Const.Value, p.Const.FloatValue, p.Const.Id)
		if err != nil {
			panic(err)
		}
	} else {

		if err := service.Const.Set(p.Const.Name, p.Const.Key, p.Const.Value, p.Const.FloatValue); err != nil {
			panic(err)
		}
	}
	// return to current page.
	// return exit.RedirectFirstValid(p.Referer, "/inventory/")
}

func (p *PreferenceIndex) CurPageSuffix() string {
	if p.Tab == "today" {
		return "/today"
	}
	return ""
}

// pager related

func (p *PreferenceIndex) UrlTemplate() string {
	return fmt.Sprintf("/admin/announcement/%s/{{Start}}/{{PageItems}}", p.Tab)
}
