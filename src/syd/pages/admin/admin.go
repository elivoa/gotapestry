package admin

import (
	"got/core"
	"got/register"
	"syd/service/productservice"
)

func Register() {}
func init() {
	register.Page(Register, &AdminIndex{})
}

//________________________________________________________________________________
//
//
type AdminIndex struct {
	core.Page
}

func (p *AdminIndex) OnRebuildProductPinyin() {
	productservice.RebuildProductPinyinCapital()
}
