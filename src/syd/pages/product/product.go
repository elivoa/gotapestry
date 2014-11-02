package product

import (
	"github.com/elivoa/got/core"
)

/* ________________________________________________________________________________
   Product Home Page
*/
type ProductIndex struct{ core.Page }

func (p *ProductIndex) Setup() (string, string) { return "redirect", "/product/list" }

// redirect
type ProductCreate struct{ core.Page }

func (p *ProductCreate) Setup() (string, string) { return "redirect", "/product/edit" }

// --------------------------------------------------------------------------------

var (
	//	listTypeLabel   = map[string]string{"customer": "客户", "factory": "厂商"}
	createEditLabel = map[string]string{"create": "新建", "edit": "编辑"}
)
