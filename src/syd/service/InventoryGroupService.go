package service

import (
	"errors"
	"fmt"
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
		if list, err := Inventory.List(parser, withs); err != nil {
			return nil, err
		} else {
			ig.Inventories = list
		}
		// invs, err := inventorydao.List(parser)
		// if err != nil {
		// 	return nil, err
		// }
	}
	// TODO this should be in list.
	if withs&WITH_PRODUCT > 0 {
		s.FillWithProducts(ig)
	}
	// if withs&WITH_STOCKS > 0 {
	// 	s.FillWithStocks(ig)
	// }
	return ig, err
}

func (s *InventoryGroupService) FillWithProducts(ig *model.InventoryGroup) error {
	var idset = map[int64]bool{}
	for _, inv := range ig.Inventories {
		idset[inv.ProductId] = true
	}
	productmap, err := Product.BatchFetchProductByIdMap(idset)
	if err != nil {
		return err
	}
	if nil != productmap && len(productmap) > 0 {
		for _, inv := range ig.Inventories {
			if product, ok := productmap[inv.ProductId]; ok {
				inv.Product = product
			}
		}
	}
	return nil
}

// func (s *InventoryGroupService) FillWithStocks(ig *model.InventoryGroup) error {
// 	var idset = map[int64]bool{}
// 	for _, inv := range ig.Inventories {
// 		idset[inv.ProductId] = true
// 	}
// 	productmap, err := Product.BatchFetchProductByIdMap(idset)
// 	if err != nil {
// 		return err
// 	}
// 	if nil != productmap && len(productmap) > 0 {
// 		for _, inv := range ig.Inventories {
// 			if product, ok := productmap[inv.ProductId]; ok {
// 				inv.Product = product
// 			}
// 		}
// 	}
// 	return nil
// }

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

// func (s *InventoryGroupService) CreateInventoryGroupOnly(ig *model.InventoryGroup) (
// 	*model.InventoryGroup, error) {
// 	return nil, nil
// }

// If InvemtoryGroup has GroupId property, Update, else Create one.
// 注意：参数传进来的invneotries的库存信息存放在Stocks里面。而数据库中读取出来的库存信息放在.Stock中。
func (s *InventoryGroupService) SaveInventoryGroupByNGLIST(ig *model.InventoryGroup) (
	*model.InventoryGroup, error) {

	fmt.Println("\n\n\n === --- >>> TODO CreateInventoryGroupByNGLIST")
	fmt.Println("inventory group is  : ", ig.Id)

	// nil check
	if nil == ig {
		return nil, errors.New("InventoryGroup can't be nil!")
	}
	if ig.ProviderId <= 0 {
		return nil, errors.New("Provider Id can't be nil!")
	}

	if ig.Id <= 0 {
		// create inventory group, make new id.
		var err error
		if ig, err = inventorygroupdao.Create(ig); err != nil {
			return nil, err
		}

		// TODO: save InventoryGroup to create GroupId.
	} else {
		// update inventory group
		fmt.Printf("update one: groupid is %d\n", ig.Id)
	}

	// assign basic properties into sub-inventories
	if nil != ig.Inventories {
		for _, i := range ig.Inventories {
			if i != nil {
				i.GroupId = ig.Id
				i.OperatorId = ig.OperatorId
				i.SendTime = ig.SendTime
				i.ReceiveTime = ig.ReceiveTime
			}
		}
	}

	// 1. load inventory from database, used to compare details.
	var (
		// details of 3 situation
		createGroup = []*model.Inventory{}
		updateGroup = []*model.Inventory{}
		deleteGroup = []*model.Inventory{}
	)

	// load existing inventories
	invs, e := inventorydao.List(db.NewQueryParser().Where(inventory.FGroupId, ig.Id))
	if e != nil {
		return nil, e
	}
	if invs == nil || len(invs) == 0 {
		// createGroup = invs // 格式不对怎办？ 这里有问题，使用Clone方法；
		for _, i := range ig.Inventories {
			if nil != i.Stocks {
				for color, sizemap := range i.Stocks {
					if sizemap != nil {
						for size, stock := range sizemap {
							if stock > 0 {
								createGroup = append(createGroup, _cloneInentory(i, color, size, stock))
							}
						}
					}
				}
			}
		}

	} else {
		var deleteWhoIsFalse = make([]bool, len(invs))
		if ig.Inventories != nil {
			for _, inv := range ig.Inventories { // 1 level loop: inv, color, size, stock
				if nil != inv && inv.Stocks != nil {
					for color, sizemap := range inv.Stocks {
						if sizemap != nil {
							for size, stock := range sizemap {

								// level 2 loop
								var find = false
								for idx2, inv2 := range invs {

									// match
									if inv2.ProductId == inv.ProductId &&
										inv2.Color == color && inv2.Size == size {

										// assign database's matched id to post inv.
										inv.Id = inv2.Id

										// if any values changes, or if quantity chagne to 0, delete it.
										if stock != inv2.Stock || inv2.Price != inv.Price ||
											inv2.Note != inv.Note {
											if stock > 0 {
												updateGroup = append(updateGroup, inv)
											}
										}
										find = true
										if stock == 0 {
											deleteWhoIsFalse[idx2] = false
										} else {
											deleteWhoIsFalse[idx2] = true
										}
										break
									}
								}
								if !find {
									newinv := _cloneInentory(inv, color, size, stock)
									createGroup = append(createGroup, newinv)
								}
							}
						}
					}
				}

			}
		}
		// who will be deleted?
		for idx, b := range deleteWhoIsFalse {
			if !b {
				deleteGroup = append(deleteGroup, invs[idx])
			}
		}
	}

	// --------------------------------------------------------------------------------
	fmt.Println(">>>> ig.Inventories and db invs:")
	// if nil != ig.Inventories {
	// 	for _, inv := range ig.Inventories {
	// 		fmt.Println("\tig.inventories: ", inv, d.Color, d.Size, " = ", d.Quantity, d.SellingPrice)
	// 	}
	// }
	if nil != invs {
		for _, inv := range invs {
			fmt.Println("\torder details: ", inv.Id, inv.Color, inv.Size, " = ", inv.Stock, inv.Price)
		}
	}
	// fmt.Println(">>>> who is false?")
	// for idx, b := range deleteWhoIsFalse {
	// 	fmt.Println("\t >> who is false: ", idx, b)
	// }
	// --------------------------------------------------------------------------------

	// final process: create, update, and delete
	fmt.Println("\n===============================================================")
	if createGroup != nil {
		for _, inventory := range createGroup {

			if _, err := inventorydao.Create(inventory); err != nil {
				return nil, err
			}
		}
	}
	if updateGroup != nil {
		for _, inventory := range updateGroup {
			if _, err := inventorydao.Update(inventory); err != nil {
				return nil, err
			}
		}
	}
	if deleteGroup != nil {
		for _, inventory := range deleteGroup {
			if _, err := inventorydao.Delete(inventory.Id); err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func _cloneInentory(inv *model.Inventory, color, size string, stock int) *model.Inventory {
	// need to clone inventory object.
	return &model.Inventory{ // clone values
		GroupId:     inv.GroupId,
		ProductId:   inv.ProductId,
		Color:       color, // 变化的值
		Size:        size,  // 变化的值
		Stock:       stock, // 变化的值
		ProviderId:  inv.ProviderId,
		OperatorId:  inv.OperatorId,
		Price:       inv.Price,
		Status:      inv.Status,
		Type:        inv.Type,
		Note:        inv.Note,
		SendTime:    inv.SendTime,
		ReceiveTime: inv.ReceiveTime,
		CreateTime:  inv.CreateTime,
		UpdateTime:  inv.UpdateTime,
	}
}

// func _processingInventoryGroupDetails(ig *model.InventoryGroup) error {
// 	return nil
// }
