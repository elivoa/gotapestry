package model

import (
	// "fmt"
	"strings"
	"time"
)

// TODO Design:model, how to split model data and fields.
type Product struct {
	Id           int    // id
	Name         string // product name
	ProductId    string // 传说中的货号
	Brand        string
	Price        float64
	Supplier     int
	FactoryPrice float64
	Stock        int    // 库存量 || not used again?
	ShelfNo      string // 货架号
	Capital      string // captical letter to quick access.
	Note         string
	CreateTime   time.Time
	UpdateTime   time.Time

	Pictures string // picture keys splited by ';' filenamne can't contain ';'

	// additional information, not in persistence
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
		CreateTime: time.Now(),
	}
}

// used for fetch;
// type ProductCSValue struct {
// 	Id        int
// 	ProductId int
// 	Color     []string
// 	Size      []string
// }

func (p *Product) TotalStock() int {
	if nil != p.Stocks && len(p.Stocks) > 0 {
		var totalstock = 0
		for _, s := range p.Stocks {
			totalstock += s
		}
		return totalstock
	}
	return 0
}

func (p *Product) ClearValues() {
	if p != nil {
		p.ClearColors()
		p.ClearSizes()
	}
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

func (p *Product) PictureKeys() []string {
	return strings.Split(p.Pictures, ";")
}
