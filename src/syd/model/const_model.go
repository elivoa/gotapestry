package model

import (
	"time"
)

type Const struct {
	Id         int64
	Name       string
	Key        string
	Value      string
	FloatValue float64
	Order      int // TODO: implement order to manage order.
	Time       time.Time
}
