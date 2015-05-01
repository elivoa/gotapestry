// Time-stamp: <[InventoryIndex.go] Elivoa @ Wednesday, 2015-04-22 16:38:38>
package inventory

import (
	"fmt"
	"github.com/elivoa/got/builtin/services"
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/db"
	"github.com/elivoa/got/route/exit"
	"github.com/elivoa/got/utils"
	"strings"
	"syd/base/inventory"
	"syd/model"
	"syd/service"
	"time"
)

type InventoryIndex struct {
	core.Page

	InventoryGroups []*model.InventoryGroup
	Tab             string `path-param:"1"`
	Current         int    `path-param:"2"` // pager: the current item. in pager.
	PageItems       int    `path-param:"3"` // pager: page size.

	// searc form or filter
	Provider int64     `query:"provider"` // filter by provider Id.
	TimeFrom time.Time `query:"from"`
	TimeTo   time.Time `query:"to"`

	// properties
	Total int // pager: total items available

	Referer string // return to this place
	// TimeZone *model.TimeZoneInfo
}

func (p *InventoryIndex) Activate() {
	// service.User.RequireRole(p.W, p.R, syd.RoleSet_Orders...)

	// process parameters
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

	// parameter time
	if !utils.IsValidTime(p.TimeFrom) {
		p.TimeFrom = time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local)
	}
	if !utils.IsValidTime(p.TimeTo) {
		p.TimeTo = time.Now()
	}

	// load inventory group
	var err error
	parser := service.InventoryGroup.EntityManager().NewQueryParser().Where()
	switch strings.ToLower(p.Tab) {
	case "all", "":
	default:
		// parser.And("status", p.Tab)
	}

	// filter by provider
	if p.Provider > 0 {
		fmt.Println("_______________")
		fmt.Println("And filter by provider_id ", p.Provider)
		parser.And(inventory.F_ProviderId, p.Provider)
	}
	// time
	if utils.IsValidTime(p.TimeFrom) && utils.IsValidTime(p.TimeTo) {
		parser.Range(inventory.FSendTime, p.TimeFrom, p.TimeTo)
	}

	parser.And(inventory.F_Type, inventory.TypeReceive)

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

func (p *InventoryIndex) OnSuccessFromSearchForm() *exit.Exit {
	// time is injected and then return linkpage.
	return exit.Redirect(p.ThisPageLink())
}

func (p *InventoryIndex) OnClearForm() *exit.Exit {
	p.TimeFrom = time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
	p.TimeTo = p.TimeFrom
	return exit.Redirect(p.ThisPageLink())
}

func (p *InventoryIndex) ThisPageLink() string {
	// 一个普通的SearchBox实现。所有东西都放到url里面。直接redirect到本页面。
	var parameters = map[string]interface{}{
		"provider": p.Provider,
		"from":     p.TimeFrom,
		"to":       p.TimeTo,
	}

	url := services.Link.GeneratePageUrlWithContextAndQueryParameters("inventory", parameters)
	return url

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
