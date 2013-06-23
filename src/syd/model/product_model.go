package model

import (
	// "fmt"
	"time"
)

type Product struct {
	Id           int    // id
	Name         string // product name
	ProductId    string // 传说中的货号
	Brand        string
	Price        float64
	Supplier     int
	FactoryPrice float64
	Stock        int // 库存量
	Note         string
	CreateTime   time.Time
	UpdateTime   time.Time
}

func NewProduct() *Product {
	// return &Product{Name: "测试:破洞猫头", ProductId: "123456", Brand: "sniddle", Price: 128.8, FactoryPrice: 10.00, Stock: 1000, Note: "这是一个测试的条目"}
	return &Product{}
}

// TODO type is enum
