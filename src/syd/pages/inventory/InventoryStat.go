package inventory

import (
	"encoding/json"
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"github.com/elivoa/gxl"
	"syd/base/person"
	"syd/model"
	"syd/service"
	"time"
)

/* ________________________________________________________________________________
   Product Craete Page
*/
type InventoryStat struct {
	core.Page

	// parameter
	Factory  int64 `query:"factory"` // factory id
	TimeFrom time.Time
	TimeTo   time.Time

	// property
	GroupId        *gxl.Int              `path-param:"1"`
	InventoryGroup *model.InventoryGroup ``

	// ** special ** for angularjs form submit. this is json.
	InventoriesJson string

	Referer string `query:"referer"` // referer page, view or list

}

func (p *InventoryStat) New() *InventoryStat {
	return &InventoryStat{}
}

func (p *InventoryStat) Setup() {
	// TODO load TimeFrom
	// if gxl.CN == p.TimeFrom

	// page values
	if p.GroupId != nil {
		var err error
		ig, err := service.InventoryGroup.GetInventoryGroup((int64)(p.GroupId.Int),
			service.WITH_INVENTORIES|service.WITH_PRODUCT|service.WITH_STOCKS)
		if err != nil {
			panic(err)
		}
		p.InventoryGroup = ig

		// parser := db.NewQueryParser().Where(inventory.FGroupId, p.GroupId.Int)
		// list, err := service.Inventory.List(parser,
		// 	service.WITH_PERSON|service.WITH_PRODUCT|service.WITH_USERS)
		// if err != nil {
		// 	panic(err)
		// }

		// construct group;
		// p.InventoryGroup = model.NewInventoryGroup(list)
	} else {
		p.InventoryGroup = model.NewInventoryGroup(nil)
	}
}

func (p *InventoryStat) Factories() []*model.Person {
	if persons, err := service.Person.GetPersons(person.TYPE_FACTORY); err != nil {
		panic(err)
	} else {
		return persons
	}
}

/** For form submit. */

func (p *InventoryStat) OnPrepareForSubmitFromInventoryForm() {
	if p.GroupId == nil { // if create
		// p.InventoryGroup = model.NewProduct()
	} else { // if edit
		// 读取了数据库的order是为了保证更新的时候不会丢失form中没有的数据；
		ig, err := service.InventoryGroup.GetInventoryGroup((int64)(p.GroupId.Int),
			0 /* service.WITH_INVENTORIES*/)
		if err != nil {
			panic(err.Error())
		}

		p.InventoryGroup = ig
	}
}

func (p *InventoryStat) OnSuccessFromInventoryForm() *exit.Exit {

	invs, err := p.unmarshalInventories(p.InventoriesJson)
	if err != nil {
		panic(err)
	}
	p.InventoryGroup.Inventories = invs

	if p.GroupId == nil { // if create
		// Auto add 2 days to ReceiveTime.
		p.InventoryGroup.ReceiveTime = p.InventoryGroup.SendTime.AddDate(0, 0, 2)
	} else { // if edit
		// TODO...
	}

	nig, err := service.InventoryGroup.SaveInventoryGroupByNGLIST(p.InventoryGroup)
	if err != nil {
		panic(err)
	}

	// return to refer first
	return exit.RedirectFirstValid(
		p.Referer,
		"/product/list",
		fmt.Sprintf("/inventory/edit/%d", nig.Id),
	)
}

// return []*model.Inventory with Stocks(temp variable) in it;
func (p *InventoryStat) unmarshalInventories(invsJson string) ([]*model.Inventory, error) {
	invs := []*model.Inventory{}
	if err := json.Unmarshal([]byte(invsJson), &invs); err == nil {
		if invs != nil {
			fmt.Println("invs is not nil")
			for idx, a := range invs {
				if a != nil {
					fmt.Println(idx, " : ", a)
					fmt.Println(" id: ", a.Id, "; productId: ", a.ProductId)
					fmt.Println("Stocks is", a.Stocks)
				}
			}
		}
	} else {
		return nil, err
	}
	return invs, nil
}
