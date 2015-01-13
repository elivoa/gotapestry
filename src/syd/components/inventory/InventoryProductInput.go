package inventory

import (
	"github.com/elivoa/got/core"
	"syd/model"
)

type InventoryProductInput struct {
	core.Component

	GroupId     int64
	Inventories []*model.Inventory // Inventory items in InventoryGroup.
	TotalPrice  float64            // all order's price
	Referer     string             // return to this place

	// TimeZone *model.TimeZoneInfo
}

func (p *InventoryProductInput) SetupRender() {
	// verify user role.
	// service.User.RequireRole(p.W, p.R, "admin") // TODO remove w, r. use service injection.
	// p.TimeZone = service.TimeZone.UserTimeZoneSafe(p.R)
	if p.Inventories == nil {
		return
	}
}
