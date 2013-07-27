package got

import (
	"fmt"
	"got/core"
	"got/register"
	"got/templates"
	"html/template"
	"syd/service/suggest"
)

func Register() {}

func init() {
	register.Page(Register, &Status{})
}

// TODO
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

func (p *Status) AfterRender() {
	fmt.Println("\n\n---------------------------\n\n")
	suggest.PrintAll()
}

func (p *Status) OnClickTemplate(name string) {
	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	fmt.Println("click template: ", name)
	t:= templates.Templates.Lookup(name)
	fmt.Println(t)
	fmt.Println(t.Delims)
}



















