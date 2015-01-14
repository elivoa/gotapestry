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
	Stock      int
	ProviderId int64 // factory person id.
	OperatorId int64 // user id.

	Status inventory.Status
	Type   inventory.Type
	Note   string

	// Operator
	SendTime    time.Time // 发货时间
	ReceiveTime time.Time // 收到货的时间
	CreateTime  time.Time // 创建收货单的时间
	UpdateTime  time.Time

	// extended:
	Product  *Product
	Provider *Person // factory
	Operator *User   // operator
}

// no need to persist to db?
type InventoryGroup struct {
	Id          int64
	Inventories []*Inventory // inventory list;

	Status inventory.Status
	Type   inventory.Type
	Note   string

	// copied from Inventory
	SendTime    time.Time // 发货时间
	ReceiveTime time.Time // 收到货的时间
	CreateTime  time.Time
	UpdateTime  time.Time
}

func NewInventoryGroup(invs []*Inventory) *InventoryGroup {
	if nil == invs || len(invs) == 0 {
		return &InventoryGroup{}
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
