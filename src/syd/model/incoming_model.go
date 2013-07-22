package model

import (
	"time"
)

type AccountIncoming struct {
	Id        int
	CustomeId int
	Incoming  float64
	Time      time.Time
}
