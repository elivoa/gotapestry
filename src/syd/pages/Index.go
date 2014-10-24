package index

import (
	"github.com/elivoa/got/core"
)

// _______________________________________________________________________________
//  ROOT Page
//
type Index struct {
	core.Page
	Title string
}

func (p *Index) SetupRender() {
	p.Title = "圣衣蝶服饰销售管理系统"
}
