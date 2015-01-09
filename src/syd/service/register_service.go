package service

type Withs int

var (
	WITH_USERS   Withs = 1 << 0
	WITH_PERSON  Withs = 1 << 1 // customer or factory
	WITH_PRODUCT Withs = 1 << 2 // customer or factory

	WITH_PRODUCT_DETAIL    Withs = 1 << 10 // color size information and inventory
	WITH_PRODUCT_INVENTORY Withs = 1 << 11
)

// 临时这样初始化service, 以后要用Inject的方式初始化这些东西；
var (
	Order     = new(OrderService)
	Account   = new(AccountService)
	Person    = new(PersonService)
	User      = new(UserService)
	Product   = new(ProductService)
	Inventory = new(InventoryService)
)
