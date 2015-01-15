package inventory

import (
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/db"
	"github.com/elivoa/gxl"
	"syd/base/inventory"
	"syd/model"
	"syd/service"
)

/* ________________________________________________________________________________
   Product Create Page
*/
type InventoryEdit struct {
	core.Page

	// field
	Title    string
	SubTitle string

	// property
	GroupId        *gxl.Int              `path-param:"1"`
	InventoryGroup *model.InventoryGroup ``

	// helper used because angularjs
	// Sizes  []*model.Object

	// Pictures []string // uploaded picture's key

	// ...
	Referer string `query:"referer"` // referer page, view or list

	// display
	StockJson string
}

// func (p *InventoryEdit) ProductJson() *model.Product {
// 	return p.InventoryGroup
// }

// init this page
func (p *InventoryEdit) New() *InventoryEdit {
	return &InventoryEdit{}
}

func (p *InventoryEdit) Setup() {
	// debug group id as 1
	// p.GroupId = gxl.NewInt(1)

	// page values
	p.Title = "create input post"
	if p.GroupId != nil {
		var err error
		parser := db.NewQueryParser().Where(inventory.FGroupId, p.GroupId.Int)
		list, err := service.Inventory.List(parser,
			service.WITH_PERSON|service.WITH_PRODUCT|service.WITH_USERS)
		if err != nil {
			panic(err)
		}

		// construct group;
		p.InventoryGroup = model.NewInventoryGroup(list)
		p.SubTitle = "编辑"
	} else {
		p.InventoryGroup = model.NewInventoryGroup(nil)
		p.SubTitle = "新建"
	}
}

func (p *InventoryEdit) InventoriesJson() []*model.Inventory {
	if p.InventoryGroup != nil {
		return p.InventoryGroup.Inventories
	}
	return nil
}

// func (p *InventoryEdit) OnPrepareForSubmitFromProductForm() {
// 	if p.Id == nil { // if create
// 		p.Product = model.NewProduct()
// 	} else {
// 		// if edit
// 		// for security reason, TODO security check here.
// 		// 读取了数据库的order是为了保证更新的时候不会丢失form中没有的数据；
// 		model, err := service.Product.GetFullProduct(p.Id.Int)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		p.Product = model
// 		// 但是这样做就必须清除form更新的时候需要删除的值，否则form提交和原有值是叠加的，会引起错误；
// 		// 这里只需要清除列表等数据，这个Order中只有Details是列表。
// 		p.Product.ClearColors()
// 		p.Product.ClearSizes()
// 		p.Product.ClearValues()
// 	}
// }

// func (p *InventoryEdit) OnSuccessFromProductForm() *exit.Exit {
// 	// clear values
// 	p.Product.ClearValues()

// 	// transfer stocks value to product.Stocks
// 	if p.Stocks != nil {
// 		p.Product.Stocks = make([]*model.ProductStockItem, len(p.Product.Colors)*len(p.Product.Sizes))

// 		i := 0
// 		for _, color := range p.Product.Colors {
// 			for _, size := range p.Product.Sizes {
// 				fmt.Println("?>>>>", i)
// 				// key := fmt.Sprintf("%v__%v", color, size)
// 				p.Product.Stocks[i] = &model.ProductStockItem{
// 					Color: color,
// 					Size:  size,
// 					Stock: p.Stocks[i],
// 				}
// 				i = i + 1
// 			}
// 		}
// 	}

// 	// transfer pictures value to pictures.
// 	if p.Pictures != nil {
// 		p.Product.Pictures = strings.Join(p.Pictures, ";")
// 	}

// 	// write to db
// 	if p.Id != nil {
// 		service.Product.UpdateProduct(p.Product)
// 	} else {
// 		service.Product.CreateProduct(p.Product)
// 	}

// 	if p.Referer == "view" {
// 		return exit.Redirect(fmt.Sprintf("/product/detail/%v", p.Product.Id))
// 	}
// 	// TODO: return to original page.
// 	return exit.Redirect("/product/list")
// }
