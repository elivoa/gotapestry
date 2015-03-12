package inventory

import (
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"syd/model"
)

type InventoryGroupList struct {
	core.Component

	InventoryGroups []*model.InventoryGroup
	TotalPrice      float64 // all order's price
	Referer         string  // return to this place
}

func (p *InventoryGroupList) SetupRender() {
	// verify user role.
	// service.User.RequireRole(p.W, p.R, "admin") // TODO remove w, r. use service injection.
	// p.TimeZone = service.TimeZone.UserTimeZoneSafe(p.R)

	if p.InventoryGroups == nil {
		return
	}
}

// ________________________________________________________________________________
// Events
//
func (p *InventoryGroupList) OnDelete(id int64) *exit.Exit {
	fmt.Printf("kkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk\n TODO: delete id , %d\n", id)
	// TODO delete inventories
	// TODO delete inventory groups.
	// if _, err := service.Inventory.DeleteInventory(id); err != nil {
	// 	panic(err)				//
	// }							//
	return exit.RedirectFirstValid(p.Referer, "/inventory")
}

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
