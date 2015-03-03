// Time-stamp: <[InventoryIndex.go] Elivoa @ Tuesday, 2015-03-03 11:09:22>
package inventory

import (
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/db"
	"strings"
	"syd/base/inventory"
	"syd/model"
	"syd/service"
)

type InventoryIndex struct {
	core.Page

	InventoryGroups []*model.InventoryGroup
	Tab             string `path-param:"1"`
	Current         int    `path-param:"2"` // pager: the current item. in pager.
	PageItems       int    `path-param:"3"` // pager: page size.

	// properties
	Total int // pager: total items available

	Referer string // return to this place

	// TimeZone *model.TimeZoneInfo
}

func (p *InventoryIndex) Activate() {
	// service.User.RequireRole(p.W, p.R, syd.RoleSet_Orders...)

	// not injected with parameters.
	if p.Tab == "" {
		p.Tab = "all" // default go in toprint
	}
}

func (p *InventoryIndex) SetupRender() {
	// verify user role.
	// service.User.RequireRole(p.W, p.R, "admin") // TODO remove w, r. use service injection.

	// fix the pagers
	if p.PageItems <= 0 {
		p.PageItems = config.LIST_PAGE_SIZE // TODO default pager number. Config this.
	}

	// load inventory group
	var err error
	parser := service.InventoryGroup.EntityManager().NewQueryParser().Where()
	switch strings.ToLower(p.Tab) {
	case "all", "":
	default:
		// parser.And("status", p.Tab)
	}
	// parser.Or("type", model.Wholesale, model.ShippingInstead) // restrict type
	parser.OrderBy(inventory.FSendTime, db.DESC)

	// get total
	p.Total, err = parser.Count()
	if err != nil {
		panic(err.Error())
	}

	// 2. get order list.
	parser.Limit(p.Current, p.PageItems) // pager
	p.InventoryGroups, err = service.InventoryGroup.List(parser, service.WITH_PERSON)
	if err != nil {
		panic(err.Error())
	}
}

// func (p *InventoryIndex) ShowProduct(r *model.Inventory) string {
// 	if nil != r.Product {
// 		return r.Product.Name
// 	} else {
// 		return strconv.FormatInt(r.ProductId, 10)
// 	}
// }

// ________________________________________________________________________________
// Events
//
// func (p *InventoryIndex) Ondelete(id int64, tab string) interface{} {
// 	if _, err := service.Inventory.DeleteInventory(id); err != nil {
// 		panic(err)
// 	}
// 	return exit.RedirectFirstValid(p.Referer, "/inventory")
// }

// func (p *InventoryList) OnMarkInUse(id int64) interface{} {
// 	token := service.User.GetLogin(p.W, p.R)
// 	if inventory, err := service.Inventory.GetInventory(id); err != nil {
// 		panic(exception.NewCoreError(err, ""))
// 	} else {
// 		if nil == inventory {
// 			panic(exception.NewCoreError(nil, "Inventory %d can't be null!", id))
// 		}
// 		if inventory.Status != model.InventoryStatus_NewPurchased {
// 			panic(exception.NewCoreError(nil, "Can't use Inventory(%d) with status %s!",
// 				id, inventory.Status))
// 		}
// 		// update inventory
// 		inventory.Status = model.InventoryStatus_InUse
// 		inventory.UseTime = time.Now()
// 		if _, err := service.Inventory.UpdateInventory(inventory); err != nil {
// 			panic(exception.NewCoreError(err, "Error updating inventory(%d)!", id))
// 		} else {
// 			// if update success, write user log.
// 			service.User.LogUserAction(token.Id, model.ACTION_MARK_INVENTORY_INUSE, id)
// 		}
// 	}
// 	return exit.RedirectFirstValid(p.Referer, "/inventory")
// }

// func (p *InventoryList) OnMarkRunout(id int64) interface{} {
// 	token := service.User.GetLogin(p.W, p.R)
// 	if inventory, err := service.Inventory.GetInventory(id); err != nil {
// 		panic(exception.NewCoreError(err, ""))
// 	} else {
// 		if nil == inventory {
// 			panic(exception.NewCoreError(nil, "Inventory %d can't be null!", id))
// 		}
// 		if inventory.Status != model.InventoryStatus_InUse {
// 			panic(exception.NewCoreError(nil, "Can't mark runout to Inventory(%d) with status %s!",
// 				id, inventory.Status))
// 		}
// 		// update inventory
// 		inventory.Status = model.InventoryStatus_RunOut
// 		inventory.RunOutTime = time.Now()
// 		if _, err := service.Inventory.UpdateInventory(inventory); err != nil {
// 			panic(exception.NewCoreError(err, "Error updating inventory(%d)!", id))
// 		} else {
// 			// if update success, write user log.
// 			service.User.LogUserAction(token.Id, model.ACTION_MARK_INVENTORY_RUNOUT, id)
// 		}
// 	}
// 	return exit.RedirectFirstValid(p.Referer, "/inventory")
// }
