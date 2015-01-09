// Time-stamp: <[inventory_data.go] Elivoa @ Sunday, 2014-11-23 14:03:33>

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
	FID      = "id"
	FGroupId = "group_id"
)
