package model

import (
	// "fmt"
	"time"
)

//
// core person model
//
type Person struct {
	Id         int    // id // TODO: change to int64
	Name       string // pesron name
	Type       string // `enum(客户Customer|厂家Factory)` // person type
	Phone      string
	City       string
	Address    string
	PostalCode int
	QQ         int
	Website    string
	Note       string

	// Customer: 存储累计欠款; Factory: TODO
	AccountBallance float64

	CreateTime time.Time
	UpdateTime time.Time

	// TODO ++
	/* favorite delivery method: TakeAway|SFExpress|物流 */
	DeliveryMethod string

	// Fax        string
}

func NewPerson() *Person {
	return &Person{Name: "", Type: "customer", Note: ""}
}

func (p *Person) Accomulated() float64 {
	return -p.AccountBallance
}

func (p *Person) IsCustomer() bool {
	return p.Type == "customer"
}

func (p *Person) IsFactory() bool {
	return p.Type == "factory"
}

//
// Advanced Wrapper
//
type Customer struct {
	Person
	// advanced properties
	// Accumulated float64 // 累计欠款 // TODO replaced by AccountBallance
}

type Producer struct {
	Person
	// advanced properties
}

// TODO type is enum

//
// Customer Special Price
//

type CustomerPrice struct {
	Id           int
	PersonId     int
	ProductId    int
	Price        float64
	CreateTime   time.Time
	LastUsedTime time.Time
}
