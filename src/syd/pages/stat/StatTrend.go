package stat

import (
	"github.com/elivoa/got/core"
)

type StatTrend struct {
	core.Page

	// Id *gxl.Int `path-param:"1"`

	Period     int `query:"period"` // Chart Time Period
	CombineDay int `query:"combineday"`
	Yearonyear int `query:"yearonyear"`

	Days int `query:"days"` // paylog days

}
