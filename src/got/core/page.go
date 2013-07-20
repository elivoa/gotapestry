package core

/*
   Interfaces of pages/component
   TODO: change interface to Pager/Protoner/Componenter
*/
type IPage interface {
	Protoner
}

/*
   Basic Page Struct
*/
type Page struct {
	Proton
}

func (p *Page) Kind() Kind {
	return PAGE
}
