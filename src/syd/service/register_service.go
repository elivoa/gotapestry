package service

// 临时这样初始化service, 以后要用Inject的方式初始化这些东西；
var (
	Order   = new(OrderService)
	Account = new(AccountService)
)
