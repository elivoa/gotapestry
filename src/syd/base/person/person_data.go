// Time-stamp: <[person_data.go] Elivoa @ Sunday, 2015-01-18 23:32:49>

package person

type Type string // TODO Change these into Person.Type
type Status int

var (
	TYPE_CUSTOMER Type = "Customer"
	TYPE_FACTORY  Type = "Factory"

	// StatusPredict Status = 1
	// StatusNormal  Status = 2
)

// Fields
var (
	f_ID      = "id"
	f_GroupId = "group_id"
)
