// Time-stamp: <[inventory_data.go] Elivoa @ Friday, 2015-12-18 12:09:27>

package inventory

type Type int
type Status int

var (
	TypeReceive     Type = 0 // 入库
	TypePlaceOrder  Type = 1 // 下单
	TypeClientOrder Type = 2 // 客户下单

	StatusPredict Status = 1
	StatusNormal  Status = 2
)

// Fields
var (
	FID          = "id"
	FGroupId     = "group_id"
	FSendTime    = "send_time"
	F_ProviderId = "provider_id"
	F_Type       = "type"
)
var (
	F_Track_ProductId = "product_id"
)
