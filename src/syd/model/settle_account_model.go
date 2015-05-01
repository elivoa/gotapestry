package model

import (
	"time"
)

type FactorySettleAccount struct {
	Id               int64
	FactoryId        int64
	GoodsDescription string
	FromTime         time.Time
	SettleTime       time.Time
	ShouldPay        float64
	Paid             float64
	Note             string
	OperatorId       int64

	// fill these
	Factory  *Person
	Operator *User
}
