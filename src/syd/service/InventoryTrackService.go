package service

import (
	"github.com/elivoa/got/db"
	itdao "syd/dal/inventorytrackdao"
	"syd/model"
)

type InventoryTrackService struct{}

func (s *InventoryTrackService) EntityManager() *db.Entity {
	return itdao.EntityManager()
}

// func (s *InventoryTrackService) Track() (*model.InventoryTrackItem, error) {
// 	if m == nil {
// 		panic("Inventory can't be null!")
// 	}
// 	return itdao.CreateInventoryTrack(m)
// }

func (s *InventoryTrackService) Create(m *model.InventoryTrackItem) (*model.InventoryTrackItem, error) {
	if m == nil {
		panic("Inventory can't be null!")
	}
	return itdao.CreateInventoryTrack(m)
}

// func (s *InventoryTrackService) UpdateInventory(m *model.Inventory) (int64, error) {
// 	if m == nil {
// 		panic("Inventory can't be null!")
// 	}
// 	// m.PrepareToUpdate() // prepare to update
// 	return inventorydao.Update(m)
// }

// func (s *InventoryTrackService) GetInventoryTrack(id int64) (*model.InventoryTrackItem, error) {
// 	return itdao.GetInventoryById(id)
// }

func (s *InventoryTrackService) Delete(id int64) (affacted int64, err error) {
	return itdao.DeleteInventoryTrack(id)
}

// With Users, product, person;
func (s *InventoryTrackService) List(parser *db.QueryParser, withs Withs) (
	[]*model.InventoryTrackItem, error) {
	tracks, err := itdao.List(parser)
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
		// if withs&WITH_STOCKS > 0 {
		// 	if err := s.FillWithStocks(inventories); err != nil {
		// 		return nil, err
		// 	}
		// }
		return tracks, nil
	}
}
