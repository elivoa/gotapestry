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

// Save all changes of account ballance, and its reason or related order id.
type AccountChangeLog struct {
	Id         int64
	CustomerId int
	Delta      float64
	Account    float64

	// Type Enum:
	// 0. ?,
	// 1. manually modified,
	// 2. order send / takeaway order create,
	// 3. batch close order, 批量结款，打钱；
	// 4. order cancel,
	Type           int
	RelatedOrderTN int64
	Reason         string
	Time           time.Time
}

// PayLog is Used to List pay log in page.
type PayLog struct {
	ID           int64
	CustomerID   int64
	CustomerName string
	Type         int
	Delta        float64
	Account      float64
	Reason       string
	Time         time.Time
}

func (p *PayLog) MinusAccount() float64 {
	if p.Account == 0 {
		return p.Account
	}
	return -p.Account
}
