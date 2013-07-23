package got

import (
	"got/core"
	"got/register"
	"got/templates"
	"html/template"
	"syd/service/suggest"
	"fmt"
)

func Register() {}

func init() {
	register.Page(Register, &Status{})
}

type Status struct {
	core.Page

	Apps       *register.AppConfigs
	Pages      *register.ProtonSegment
	Components *register.ProtonSegment
	Tpls       []*template.Template
}

func (p *Status) SetupRender() {
	p.Tpls = templates.Templates.Templates()
	p.Apps = register.Apps
	p.Pages = &register.Pages
}

func (p *Status) AfterRender(){
	fmt.Println("\n\n---------------------------\n\n")
	suggest.PrintAll()
}
