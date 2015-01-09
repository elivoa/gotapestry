package inventory

import (
	"bytes"
	"fmt"
	"github.com/elivoa/got/builtin/services"
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"strings"
	"syd/model"
	"syd/service"
	"time"
)

/* ________________________________________________________________________________
The Order List page
*/
type InventoryList struct {
	core.Page

	// parameters
	Inventories []*model.Inventory
	Tab         string `path-param:"1"` // tab: today, yestoday, all
	Current     int    `path-param:"2"` // pager: the current item. in pager.
	PageItems   int    `path-param:"3"` // pager: page size.

	FilterStores []string `query:"stores"` // used for submit

	// properties
	Total int // pager: total items available

}

func (p *InventoryList) Activate() {
	// service.User.RequireRole(p.W, p.R, carfilm.RoleSet_Inventory...)
}

func (p *InventoryList) SetupRender() {
	// fix the pagers
	if p.PageItems <= 0 {
		p.PageItems = config.LIST_PAGE_SIZE
	}

	var err error
	var parser = service.Inventory.EntityManager().NewQueryParser()
	if strings.ToLower(p.Tab) == "today" {
		now := time.Now().UTC()
		start := now.Truncate(time.Hour * 24)
		end := now.AddDate(0, 0, 1).Truncate(time.Hour * 24)
		parser.Where().Range("create_time", start, end)
	} else if strings.ToLower(p.Tab) == "" { // inventory, use status inventory or newpurchased
		// parser.Where("status", model.InventoryStatus_InUse).Or("status", model.InventoryStatus_NewPurchased)
	} else {
		// all, leave it there.
		parser.Where()
	}

	// add store
	if p.FilterStores != nil && len(p.FilterStores) > 0 {
		var cond bytes.Buffer
		var values = []interface{}{}
		cond.WriteRune('(')
		for idx, store := range p.FilterStores {
			if idx > 0 {
				cond.WriteString(" or ")
			}
			cond.WriteString("store=?")
			values = append(values, store)
			// cond.WriteString("store=\"")
			// cond.WriteString(store)
			// cond.WriteString("\"")
		}
		cond.WriteRune(')')
		parser.AndRaw(cond.String(), values...)
	}

	// get total
	p.Total, err = parser.Count()
	if err != nil {
		panic(err.Error())
	}

	// 2. get order list.
	parser.Limit(p.Current, p.PageItems)
	p.Inventories, err = service.Inventory.List(parser,
		service.WITH_PERSON|service.WITH_PRODUCT|service.WITH_USERS)
	if err != nil {
		panic(err.Error())
	}
}

func (p *InventoryList) TabStyle(tab string) string {
	if strings.ToLower(p.Tab) == strings.ToLower(tab) {
		return "cur"
	}
	return ""
}

func (p *InventoryList) CurPageSuffix() string {
	if p.Tab == "history" {
		return "/history"
	}
	return ""
}

// pager related

func (p *InventoryList) UrlTemplate() string {
	return fmt.Sprintf("/inventory/%s/{{Start}}/{{PageItems}}", p.Tab)
}

func (p *InventoryList) OnSuccessFromStoreSelectorForm() *exit.Exit {
	var stores bytes.Buffer
	if nil != p.FilterStores {
		for idx, v := range p.FilterStores {
			if idx > 0 {
				stores.WriteRune(',')
			}
			stores.WriteString(v)
		}
	}
	params := map[string]interface{}{
		"stores": stores.String(),
	}
	url := services.Link.GeneratePageUrlWithContextAndQueryParameters("inventory",
		params, p.Tab, p.Current, p.PageItems)
	fmt.Println("redirect to : ", url)
	return exit.Redirect(url)
}
