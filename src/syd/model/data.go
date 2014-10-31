package model

import (
	"time"
)

type SearchCondition struct {
	StartTime time.Time
	EndTime   time.Time
}

// if endtime not include time, make this time the last minute of the day. to optimize search
// func (s *SearchCondition) OptimizeBeforeSearch() {
// 	if h, m, sec := s.EndTime.Clock(); h+m+sec == 0 {
// 		s.EndTime.After(time.Date(0, 0, 0, 23, 59, 59, 9999, time.UTC))
// 	}
// }
