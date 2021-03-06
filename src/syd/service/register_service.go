package service

import (
	"github.com/elivoa/got/logs"
)

type Withs int

// try service API design.
var (
	WITH_ALL     Withs = 999999999
	WITH_NONE    Withs = 0 // TODO replace with  WITH_NOTHING
	WITH_NOTHING Withs = 0

	WITH_USERS   Withs = 1 << 0
	WITH_PERSON  Withs = 1 << 1 // customer or factory
	WITH_PRODUCT Withs = 1 << 2 // customer or factory

	WITH_PRODUCT_DETAIL    Withs = 1 << 10 // color size information and inventory
	WITH_PRODUCT_INVENTORY Withs = 1 << 11 // 是否返回产品的库存信息

	WITH_INVENTORIES Withs = 1 << 15 // 包含Inventories列表
	WITH_STOCKS      Withs = 1 << 16 // 包含Inventlry的库存数量.
)

// 临时这样初始化service, 以后要用Inject的方式初始化这些东西；
var (

	// Fundamental Serivces
	Const = NewCosntService()
	User  = &UserService{logs: logs.Get("SERVICE:USER:LoginCheck")}

	// basic logic services
	Order          = new(OrderService)
	OrderReturns   = new(OrderReturnsService)
	Account        = new(AccountService)
	Person         = new(PersonService)
	Product        = new(ProductService)
	Inventory      = new(InventoryService)      // 入库
	InventoryGroup = new(InventoryGroupService) // 入库组
	InventoryTrack = new(InventoryTrackService) // Inventory Track
	Stock          = new(StockService)          // 库存数量
	Stat           = new(StatService)           //

	// Extend Services
	FactorySettleAccount = new(FactorySettleAccountService) //
	SendNewProduct       = new(SendNewProductService)
)
