package service

import (
	"fmt"
	"github.com/elivoa/got/db"
	"syd/dal/inventorydao"
	"syd/model"
	"time"
)

type InventoryService struct{}

func (s *InventoryService) EntityManager() *db.Entity {
	return inventorydao.EntityManager()
}

func (s *InventoryService) CreateInventory(m *model.Inventory) (*model.Inventory, error) {
	if m == nil {
		panic("Inventory can't be null!")
	}
	m.CreateTime = time.Now()
	// m.PrepareToSave() // prepare to save, create time...
	return inventorydao.Create(m)
}

func (s *InventoryService) UpdateInventory(m *model.Inventory) (int64, error) {
	if m == nil {
		panic("Inventory can't be null!")
	}
	// m.PrepareToUpdate() // prepare to update
	return inventorydao.Update(m)
}

func (s *InventoryService) GetInventory(id int64) (*model.Inventory, error) {
	return inventorydao.GetInventoryById(id)
}

// Get inventory by serial, and load inner object product.
// If product not found, return inventory with empty product
// func (s *InventoryService) GetInventoryBySerial(serial string) (*model.Inventory, error) {
// 	if inv, err := inventorydao.Get("serialno", serial); err != nil {
// 		return nil, err
// 	} else if inv != nil {
// 		if p, err := productdao.Get(int(inv.ProductId)); err != nil {
// 			log.Printf("Error when load Product of inventory. Error is: %v", err)
// 			// return nil, err
// 		} else {
// 			inv.Product = p
// 		}
// 		if inv.Product == nil {
// 			inv.Product = &model.Product{Name: "<!ProductNotFound>"}
// 		}
// 		return inv, nil
// 	}
// 	return nil, nil
// }

func (s *InventoryService) DeleteInventory(id int64) (affacted int64, err error) {
	return inventorydao.Delete(id)
}

// With Users, product, person;
func (s *InventoryService) List(parser *db.QueryParser, withs Withs) ([]*model.Inventory, error) {
	inventories, err := inventorydao.List(parser)
	if err != nil {
		return nil, err
	} else {
		if withs&WITH_USERS > 0 {
			if err := s.FillWithUser(inventories); err != nil {
				return nil, err
			}
		}
		if withs&WITH_PRODUCT > 0 {
			if err := s.FillWithProducts(inventories); err != nil {
				return nil, err
			}
		}
		if withs&WITH_PERSON > 0 {
			if err := s.FillWithPersons(inventories); err != nil {
				return nil, err
			}
		}
		if withs&WITH_STOCKS > 0 {
			if err := s.FillWithStocks(inventories); err != nil {
				return nil, err
			}
		}
		return inventories, nil
	}
}

func (s *InventoryService) FillWithStocks(models []*model.Inventory) error {
	var idset = map[int64]bool{}
	for _, m := range models {
		idset[m.ProductId] = true
	}
	fmt.Println(">> ", idset)
	// TOOD fetch with it;
	if stocks, err := inventorydao.GetAllStocksByIdSet(idset); err != nil {
		return err
	} else {
		fmt.Println("", stocks)

		if nil != stocks {
			for _, inv := range models {
				if colors, ok := stocks[inv.ProductId]; ok && nil != colors {
					if sizes, ok := colors[inv.Color]; ok && nil != sizes {
						if stock, ok := sizes[inv.Size]; ok && stock > 0 {
							inv.LeftStock = stock
						} else {
							inv.LeftStock = 0
						}
					}
				}
			}
		}
	}
	return nil
}

// func (s *InventoryService) SearchInventoryInUseByPattern(pattern string) ([]*model.Inventory, error) {
// 	return inventorydao.SearchInventoryInUseByPattern(pattern)
// }

// func (s *InventoryService) GetInUseInventoryOptions(store string) ([][]string, error) {
// 	parser := inventorydao.EntityManager().NewQueryParser()
// 	parser.Where("status", model.InventoryStatus_InUse)
// 	if store != "" {
// 		parser.And("store", store)
// 	}
// 	parser.OrderBy("product_id", "asc").Limit(100)
// 	list, err := inventorydao.ListInventory(parser)
// 	if err != nil {
// 		return [][]string{[]string{"ERROR", "ERROR! ERROR Get Options!"}}, nil
// 	}

// 	// fill
// 	if err := s.FillWithProducts(list); err != nil {
// 		return nil, err
// 	}

// 	options := [][]string{}
// 	for _, model := range list {
// 		var display string
// 		if model.Product != nil {
// 			display = fmt.Sprintf("%s - %s,%s,%s", model.SerialNo,
// 				model.Product.Name, model.Product.Type, model.Product.Property)
// 		} else {
// 			display = fmt.Sprintf("%s - <nil>", model.SerialNo)
// 		}
// 		options = append(options, []string{model.SerialNo, display})
// 	}
// 	return options, nil
// }

// ------------- here

// list top orders with limit MAX_LIST_ITEMS(default 50)
// func (s *InventoryService) ListInventory(status string) ([]*model.Inventory, error) {
// 	replenishments, err := inventorydao.ListInventoryPager(0, carfilm.CONST_DB_DEFAULT_MAX_ITEMS)
// 	if err != nil {
// 		return nil, err
// 	} else {
// 		if err := s.FillWithUser(replenishments); err != nil {
// 			return nil, err
// 		}
// 		return replenishments, nil
// 	}
// }

// func (s *InventoryService) ListInventoryPagerWithUsers(limit int, n int) ([]*model.Inventory, error) {
// 	if replenishments, err := inventorydao.ListInventoryPager(limit, n); err != nil {
// 		return nil, err
// 	} else {
// 		if err := s.FillWithUser(replenishments); err != nil {
// 			return nil, err
// 		}
// 		return replenishments, nil
// 	}
// }

// list is passed by pointer.
func (s *InventoryService) FillWithUser(models []*model.Inventory) error {
	var idset = map[int64]bool{}
	for _, m := range models {
		idset[m.OperatorId] = true
	}
	usermap, err := User.BatchFetchUsersByIdMap(idset)
	if err != nil {
		return err
	}
	if nil != usermap {
		for _, m := range models {
			if user, ok := usermap[m.OperatorId]; ok {
				m.Operator = user
			}
		}
	}
	return nil
}

func (s *InventoryService) FillWithProducts(models []*model.Inventory) error {
	var idset = map[int64]bool{}
	for _, m := range models {
		idset[m.ProductId] = true
	}

	productmap, err := Product.BatchFetchProductByIdMap(idset)
	if err != nil {
		return err
	}
	if nil != productmap {
		for _, m := range models {
			if product, ok := productmap[m.ProductId]; ok {
				m.Product = product
			}
		}
	}
	return nil
}

func (s *InventoryService) FillWithPersons(models []*model.Inventory) error {
	var idset = map[int64]bool{}
	for _, m := range models {
		idset[m.ProviderId] = true
	}
	personmap, err := Person.BatchFetchPersonByIdMap(idset)
	if err != nil {
		return err
	}
	if nil != personmap {
		for _, m := range models {
			if provider, ok := personmap[m.ProviderId]; ok {
				m.Provider = provider
			}
		}
	}
	return nil
}
