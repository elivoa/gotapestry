package service

import (
	"github.com/elivoa/got/db"
	"syd/base/inventory"
	"syd/dal/inventorydao"
	"syd/dal/inventorygroupdao"
	"syd/model"
)

type InventoryGroupService struct{}

func (s *InventoryGroupService) EntityManager() *db.Entity {
	return inventorygroupdao.EntityManager()
}

// Get InventoryGroup
func (s *InventoryGroupService) GetInventoryGroup(id int64, withs Withs) (*model.InventoryGroup, error) {
	ig, err := inventorygroupdao.GetInventoryGroupById(id)
	if err != nil {
		return nil, err
	}
	// load inventories.
	if withs&WITH_INVENTORIES > 0 {
		parser := db.NewQueryParser().Where(inventory.FGroupId, id)
		invs, err := inventorydao.List(parser)
		if err != nil {
			return nil, err
		}
		ig.Inventories = invs
	}
	return ig, err
}

// With Users, product, person;
func (s *InventoryGroupService) List(parser *db.QueryParser, withs Withs) ([]*model.InventoryGroup, error) {
	igs, err := inventorygroupdao.List(parser)
	if err != nil {
		return nil, err
	} else {
		// if withs&WITH_USERS > 0 {
		// 	if err := s.FillWithUser(inventories); err != nil {
		// 		return nil, err
		// 	}
		// }
		// if withs&WITH_PRODUCT > 0 {
		// 	if err := s.FillWithProducts(inventories); err != nil {
		// 		return nil, err
		// 	}
		// }
		// if withs&WITH_PERSON > 0 {
		// 	if err := s.FillWithPersons(inventories); err != nil {
		// 		return nil, err
		// 	}
		// }
		return igs, nil
	}
}
