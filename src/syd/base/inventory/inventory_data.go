// Time-stamp: <[inventory_data.go] Elivoa @ Tuesday, 2015-03-03 11:04:48>

package inventory

type Type int
type Status int

var (
	TypePredict Type = 1
	TypeNormal  Type = 2

	StatusPredict Status = 1
	StatusNormal  Status = 2
)

// Fields
var (
	FID       = "id"
	FGroupId  = "group_id"
	FSendTime = "send_time"
)
