package model

import (
	// "fmt"
	"time"
)

//
// core person model
//
type Person struct {
	Id         int    // id
	Name       string // pesron name
	Type       string "enum(客户|厂家)" // person type
	Phone      string
	City       string
	Address    string
	PostalCode int
	QQ         int
	Website    string
	Note       string

	CreateTime time.Time
	UpdateTime time.Time

	// TODO ++
	/* favorite delivery method: TakeAway|SFExpress|物流 */
	DeliveryMethod string
}

func NewPerson() *Person {
	return &Person{Name: "", Type: "customer", Note: ""}
}

//
// Advanced Wrapper
//
type Customer struct {
	Person
	// advanced properties
	Accumulated float64 // 累计欠款
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
