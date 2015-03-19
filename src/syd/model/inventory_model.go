package model

import (
	"syd/base/inventory"
	"time"
)

type Inventory struct {
	Id         int64
	GroupId    int64
	ProductId  int64
	Color      string
	Size       string
	Stock      int   // should be named Quantity. Stock should be leftStock.
	ProviderId int64 // factory person id.
	OperatorId int64 // user id.

	Price string // 价格

	Status inventory.Status
	Type   inventory.Type
	Note   string

	// Operator
	SendTime    time.Time // 发货时间
	ReceiveTime time.Time // 收到货的时间
	CreateTime  time.Time // 创建收货单的时间
	UpdateTime  time.Time

	// extended:
	LeftStock int // 剩余库存量; 可以被填充.
	Product   *Product
	Provider  *Person // factory
	Operator  *User   // operator

	// Used in json, when unmarshal from client json;
	Stocks map[string]map[string]int // color, size, stock
}

// no need to persist to db?
type InventoryGroup struct {
	Id          int64
	Inventories []*Inventory // inventory list;

	Status inventory.Status
	Type   inventory.Type
	Note   string

	ProviderId int64 // factory person id.
	OperatorId int64 // user id.

	// some statistics for view only, add these fields to db.
	ProductKindSummary string
	TotalQuantity      int

	// copied from Inventory
	SendTime    time.Time // 发货时间
	ReceiveTime time.Time // 收到货的时间
	CreateTime  time.Time
	UpdateTime  time.Time

	// extended:
	Provider *Person // factory
	Operator *User   // operator
}

// parameter here is useless, change this into empty parameter methods;
func NewInventoryGroup(invs []*Inventory) *InventoryGroup {
	if nil == invs || len(invs) == 0 {
		return &InventoryGroup{
			SendTime:    time.Now().AddDate(0, 0, -1),
			ReceiveTime: time.Now().AddDate(0, 0, 1),
		}
	}
	first := invs[0]
	if nil == first {
		return &InventoryGroup{}
	}
	return &InventoryGroup{
		Id:          first.GroupId,
		Inventories: invs,
		Status:      first.Status,
		Type:        first.Type,
		Note:        first.Note,
		SendTime:    first.SendTime,
		ReceiveTime: first.ReceiveTime,
		CreateTime:  first.CreateTime,
		UpdateTime:  first.UpdateTime,
	}
}

// ----------------------------------------------------------------------------------------------------
// Inventory Change log
//

// InventoryTrackItem Stores all stock changes into database, used when track error.
type InventoryTrackItem struct {
	Id            int64
	ProductId     int64
	Color         string
	Size          string
	StockChagneTo int
	OldStock      int
	Delta         int
	UserId        int64
	Reason        string
	Context       string
	Time          time.Time
}
