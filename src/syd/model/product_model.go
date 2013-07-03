package model

import (
	// "fmt"
	"strings"
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
	Stock        int // 库存量 || not used again?
	Note         string
	CreateTime   time.Time
	UpdateTime   time.Time

	// additional information
	Colors     []string // these two information stores in product_properties table.
	Sizes      []string
	Properties map[string][]string // other properties // TODO

	// stock information. format: map[color__size]nstock
	// special values in stock table
	//   stock = -1 means this pair of combination doesn't exist.
	//   stock = -2 means the pair is deleted.(may be price is available)
	Stocks map[string]int
}

// Create default empty Product
func NewProduct() *Product {
	return &Product{
		Colors: []string{"", "", ""},
		Sizes:  []string{"", "", ""},
		// Stocks: map[string]int{},
	}
}

func (p *Product) ClearValues() {
	p.ClearColors()
	p.ClearSizes()
}

func (p *Product) ClearColors() {
	newColors := []string{}
	if p.Colors != nil {
		hasValue := false
		for _, color := range p.Colors {
			color = strings.Trim(color, " ")
			if color != "" {
				hasValue = true
				newColors = append(newColors, color)
			}
		}
		if !hasValue {
			newColors = append(newColors, "默认颜色")
		}
	}
	p.Colors = newColors
}

func (p *Product) ClearSizes() {
	newSizes := []string{}
	if p.Sizes != nil {
		hasValue := false
		for _, size := range p.Sizes {
			size = strings.Trim(size, " ")
			if size != "" {
				hasValue = true
				newSizes = append(newSizes, size)
			}
		}
		if !hasValue {
			newSizes = append(newSizes, "均码")
		}
	}
	p.Sizes = newSizes
}
