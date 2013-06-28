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

	// additional information
	Colors     []string // these two information stores in product_properties table.
	Sizes      []string
	Properties map[string][]string // other properties // TODO
}

func NewProduct() *Product {
	return &Product{
		Colors: []string{"请输入颜色"},
		Sizes:  []string{"请输入尺码"},
	}
}
