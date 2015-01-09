package inventory

import (
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"strconv"
	"syd/model"
	"syd/service"
)

type InventoryProductSelector struct {
	core.Component

	Inventories []*model.Inventory
	TotalPrice  float64 // all order's price
	Referer     string  // return to this place

	// TimeZone *model.TimeZoneInfo
}

// func (p *InventoryList) New() *InventoryList {
// 	return &InventoryList{
// 		inventoryService: service.Inventory,
// 	}
// }

func (p *InventoryProductSelector) SetupRender() {
	// verify user role.
	// service.User.RequireRole(p.W, p.R, "admin") // TODO remove w, r. use service injection.
	// p.TimeZone = service.TimeZone.UserTimeZoneSafe(p.R)
	if p.Inventories == nil {
		return
	}
}

func (p *InventoryProductSelector) ShowProduct(r *model.Inventory) string {
	if nil != r.Product {
		return r.Product.Name
	} else {
		return strconv.FormatInt(r.ProductId, 10)
	}
}

// ________________________________________________________________________________
// Events
//
func (p *InventoryProductSelector) Ondelete(id int64, tab string) interface{} {
	if _, err := service.Inventory.DeleteInventory(id); err != nil {
		panic(err)
	}
	return exit.RedirectFirstValid(p.Referer, "/inventory")
}
