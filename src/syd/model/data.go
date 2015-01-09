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

// In angularjs, simple array can't bind. Should use like this: [{value:4},{value:5},...]
// This model is used to change any array into object array;
type Object struct {
	Value interface{}
}

func NewObject(value interface{}) *Object {
	return &Object{Value: value}
}

func NewObjectArray(n int, fillValue interface{}) []*Object {
	if n <= 0 {
		return []*Object{}
	}
	array := make([]*Object, n)
	for i := 0; i < n; i++ {
		array[i] = &Object{Value: fillValue}
	}
	return array
}
