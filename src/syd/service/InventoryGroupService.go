package service

import (
	"bytes"
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
	// TODO-001 有bug， ig可能是nil。
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
		if withs&WITH_PERSON > 0 {
			if err := s.FillWithPersons(igs); err != nil {
				return nil, err
			}
		}
		return igs, nil
	}
}

func (s *InventoryGroupService) FillWithPersons(models []*model.InventoryGroup) error {
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

// If InvemtoryGroup has GroupId property, Update, else Create one.
// 注意：参数传进来的invneotries的库存信息存放在Stocks里面。而数据库中读取出来的库存信息放在.Stock中。
func (s *InventoryGroupService) SaveInventoryGroupByNGLIST(ig *model.InventoryGroup) (
	*model.InventoryGroup, error) {

	// var need_update_provider = false // if provider, sendTime, receiveTime changed.

	// nil check
	if nil == ig {
		return nil, errors.New("InventoryGroup can't be nil!")
	}
	if ig.ProviderId <= 0 {
		// return nil, errors.New("Provider Id can't be nil!")
		fmt.Printf("Provider Id can't be nil!\n")
	}

	// update summary and total quantity.
	ig.Summary, ig.TotalQuantity = _makeSummary(ig.Inventories)

	if ig.Id <= 0 {
		// create inventory group, make new id.
		var err error
		if ig, err = inventorygroupdao.Create(ig); err != nil {
			return nil, err
		}

		// TODO: save InventoryGroup to create GroupId.
	} else {
		// update inventory group
		if _, err := inventorygroupdao.Update(ig); err != nil {
			return nil, err
		}
	}

	// assign basic properties into sub-inventories. Make summary and total quantity.
	if nil != ig.Inventories {
		for _, i := range ig.Inventories {
			if i != nil {
				i.GroupId = ig.Id
				i.OperatorId = ig.OperatorId
				i.SendTime = ig.SendTime
				i.ReceiveTime = ig.ReceiveTime
				i.Type = ig.Type
			}
		}
	}

	// 1. load inventory from database, used to compare details.
	var (
		// details of 3 situation
		createGroup = []*model.Inventory{}
		updateGroup = []*model.Inventory{}
		deleteGroup = []*model.Inventory{}
		oldStocks   = map[int64]int{} // invId->stock
	)

	// load existing inventories. invs is inventories in db.
	invs, e := inventorydao.List(db.NewQueryParser().Where(inventory.FGroupId, ig.Id))
	if e != nil {
		return nil, e
	}
	if invs == nil || len(invs) == 0 {
		// DB中无信息，全部存入createGroup！
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

			// Loop Level 1: inv, color, size, stock
			for _, inv := range ig.Inventories {
				if nil != inv && inv.Stocks != nil {
					for color, sizemap := range inv.Stocks {
						if sizemap != nil {
							for size, stock := range sizemap {

								// Loop Level 2: inv2 is inventory in db.
								var find = false
								for idx2, inv2 := range invs {

									// match
									if inv2.ProductId == inv.ProductId &&
										inv2.Color == color && inv2.Size == size {

										// assign database's matched id to post inv.
										currentSubInv := _cloneInentory(inv, color, size, stock)
										currentSubInv.Id = inv2.Id

										// if any values changes(sotck, price or note),
										// or if quantity/stock chagne to 0, delete it.
										if stock != inv2.Stock || inv2.Price != currentSubInv.Price ||
											inv2.Note != currentSubInv.Note {
											//这里会有bug，导致丢失。修改信息；
											if stock > 0 {
												updateGroup = append(updateGroup, currentSubInv)
												oldStocks[currentSubInv.Id] = inv2.Stock
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
				inv := invs[idx]
				deleteGroup = append(deleteGroup, inv)
				oldStocks[inv.Id] = inv.Stock // cache old values
			}
		}
	}

	// --------------------------------------------------------------------------------
	// fmt.Println(">>>> ig.Inventories and db invs:")
	// if nil != ig.Inventories {
	// 	for _, inv := range ig.Inventories {
	// 		fmt.Println("\tig.inventories: ", inv, d.Color, d.Size, " = ", d.Quantity, d.SellingPrice)
	// 	}
	// }
	// if nil != invs {
	// 	for _, inv := range invs {
	// 		fmt.Println("\torder details: ", inv.Id, inv.Color, inv.Size, " = ", inv.Stock, inv.Price)
	// 	}
	// }
	// debug upgrade group::
	// fmt.Println(">>>> ig.Inventories> upgrade groups: ")
	// if updateGroup != nil {
	// 	for _, d := range updateGroup {
	// 		fmt.Println("\tig.inventories: ", d)
	// 		fmt.Println("\t\t old stock: ", oldStocks[d.ProductId])
	// 		fmt.Println("\t\t new stock: ", d.Stock)
	// 		//d, d.Color, d.Size, " = ", d.Quantity, d.SellingPrice
	// 	}
	// }
	// fmt.Println(">>>> who is false?")
	// for idx, b := range deleteWhoIsFalse {
	// 	fmt.Println("\t >> who is false: ", idx, b)
	// }
	// --------------------------------------------------------------------------------

	// final process: create, update, and delete
	if createGroup != nil {
		for _, inv := range createGroup {
			// create inventory item;
			inv.ProviderId = ig.ProviderId // force set providerId.
			if _, err := inventorydao.Create(inv); err != nil {
				return nil, err
			}
			// increase left stock.
			if oldStock, newStock, err := Stock.UpdateStockDelta(
				inv.ProductId, inv.Color, inv.Size, inv.Stock); err != nil {
				return nil, err
			} else {
				// TODO create new Track method to create track.
				if _, err := InventoryTrack.Create(
					&model.InventoryTrackItem{
						ProductId:     inv.ProductId,
						Color:         inv.Color,
						Size:          inv.Size,
						StockChagneTo: newStock,  // unknown
						OldStock:      oldStock,  // unknown
						Delta:         inv.Stock, // delta is not calculated.
						UserId:        0,         // TODO
						Reason:        "新增入库",
						Context:       fmt.Sprint(ig.Id),
					},
				); err != nil {
					return nil, err
				}
			}
		}
	}

	fmt.Println(">>>> ig.Inventories> upgrade groups: ") // print debug info.
	if updateGroup != nil {
		for _, inv := range updateGroup {
			// update
			inv.ProviderId = ig.ProviderId // force set providerId.
			if _, err := inventorydao.Update(inv); err != nil {
				return nil, err
			}
			// increase stock = leftstock - inv.Stock// TODO... ehrererere
			oldstock := oldStocks[inv.Id]
			delta := inv.Stock - oldstock
			// modify left stock
			if oldStock, newStock, err := Stock.UpdateStockDelta(
				inv.ProductId, inv.Color, inv.Size, delta); err != nil {
				return nil, err
			} else {
				// TODO create new Track method to create track.
				if _, err := InventoryTrack.Create(
					&model.InventoryTrackItem{
						ProductId:     inv.ProductId,
						Color:         inv.Color,
						Size:          inv.Size,
						StockChagneTo: newStock, // unknown
						OldStock:      oldStock, // unknown
						Delta:         delta,    // delta is not calculated.
						UserId:        0,        // TODO
						Reason:        "修改入库",
						Context:       fmt.Sprint(ig.Id),
					},
				); err != nil {
					return nil, err
				}
			}
		}
	}

	fmt.Println(">>>> ig.Inventories> delete group: ") // print debug info.

	if deleteGroup != nil {
		for _, inv := range deleteGroup {
			// delete inventory
			if _, err := inventorydao.Delete(inv.Id); err != nil {
				return nil, err
			}
			// decrease the stocks;
			oldstock := oldStocks[inv.Id]

			// modify left stock
			if oldStock, newStock, err := Stock.UpdateStockDelta(
				inv.ProductId, inv.Color, inv.Size, -oldstock); err != nil {
				return nil, err
			} else {
				// TODO create new Track method to create track.
				if _, err := InventoryTrack.Create(
					&model.InventoryTrackItem{
						ProductId:     inv.ProductId,
						Color:         inv.Color,
						Size:          inv.Size,
						StockChagneTo: newStock,  // unknown
						OldStock:      oldStock,  // unknown
						Delta:         -oldstock, // delta is not calculated.
						UserId:        0,         // TODO
						Reason:        "删除入库",
						Context:       fmt.Sprint(ig.Id),
					},
				); err != nil {
					return nil, err
				}
			}

		}
	}

	fmt.Println(">>>> ig.Inventories> UpdateAllInventoryItems: ") // print debug info.

	// Add by gb @ 2016-04-29: set all sub-inventory item's send_time, update_time, and factory_id.
	if err := inventorydao.UpdateAllInventoryItems(ig); err != nil {
		panic(err)
	}
	fmt.Println(">>>> ig.Inventories> All Done; ") // print debug info.

	return ig, nil
}

func _makeSummary(inventories []*model.Inventory) (summary string, totalQuantity int) {
	var buff bytes.Buffer
	if nil != inventories {
		for idx, inv := range inventories {
			if nil == inv { // bogao@2016-04-29: fix nil pointer bug.
				continue
			}

			if idx > 0 {
				buff.WriteString(", ")
			}

			if inv.Product != nil {
				buff.WriteString(inv.Product.Name)
			} else {
				buff.WriteString(fmt.Sprintf("[%d]", inv.ProductId))
			}
			// buff.WriteTo(w io.Writer)
			if nil != inv && inv.Stocks != nil {
				for color, sizemap := range inv.Stocks {
					if sizemap != nil {
						for size, stock := range sizemap {
							totalQuantity += stock
							if color == size {
								// remove unused error.
							}
						}
					}
				}
			}
		}
	}
	return buff.String(), totalQuantity
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
