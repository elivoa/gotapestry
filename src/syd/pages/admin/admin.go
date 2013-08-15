package admin

import (
	"got/core"
	"syd/service/productservice"
)

//________________________________________________________________________________
//
//
type AdminIndex struct {
	core.Page
}

func (p *AdminIndex) OnRebuildProductPinyin() {
	productservice.RebuildProductPinyinCapital()
}
