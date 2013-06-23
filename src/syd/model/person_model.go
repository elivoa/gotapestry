package model

import (
	// "fmt"
	"time"
)

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

// TODO type is enum
