package index

import (
	"got/core"
	"got/register"
)

func Register() {}
func init() {
	// register.Page(Register, &PersonIndex{}, &PersonList{}, &PersonEdit{})
	register.Page(Register, &Index{})
}

/*_______________________________________________________________________________
  ROOT Page
*/
type Index struct {
	core.Page
	Title string
}

func (p *Index) SetupRender() {
	p.Title = "圣衣蝶服饰销售管理系统"
}
