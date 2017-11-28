package inventory

import (
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"strconv"
	"syd/model"
	"syd/service"
)

type PersonSelector struct {
	core.Component

	GroupId     int64
	Inventories []*model.Inventory // Inventory items in InventoryGroup.
	TotalPrice  float64            // all order's price
	Referer     string             // return to this place

	// TimeZone *model.TimeZoneInfo
}

// func (p *InventoryList) New() *InventoryList {
// 	return &InventoryList{
// 		inventoryService: service.Inventory,
// 	}
// }

func (p *PersonSelector) SetupRender() {
	// verify user role.
	// service.User.RequireRole(p.W, p.R, "admin") // TODO remove w, r. use service injection.
	// p.TimeZone = service.TimeZone.UserTimeZoneSafe(p.R)
	if p.Inventories == nil {
		return
	}
}

func (p *PersonSelector) ShowProduct(r *model.Inventory) string {
	if nil != r.Product {
		return r.Product.Name
	} else {
		return strconv.FormatInt(r.ProductId, 10)
	}
}

// ________________________________________________________________________________
// Events
//
func (p *PersonSelector) Ondelete(id int64, tab string) interface{} {
	if _, err := service.Inventory.DeleteInventory(id); err != nil {
		panic(err)
	}
	return exit.RedirectFirstValid(p.Referer, "/inventory")
}
