package model

import (
	"fmt"
	"sort"
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

// 二维表，
// x: date - int64
// x:productId - int64
// value: - int
//    0: 没发
//    -1: 被标记为不可以发.
//    20150908: 发样衣的日期.
// Note: 只取1自然年内数据.

type ProductSalesTable struct {
	data         map[string]int
	dateMap      map[string]int
	productMap   map[int64]int
	DateLabel    []string
	ProductLabel []int64
}

func NewProductSalesTable() *ProductSalesTable {
	return &ProductSalesTable{
		data:         map[string]int{},
		dateMap:      map[string]int{},
		productMap:   map[int64]int{},
		DateLabel:    []string{},
		ProductLabel: []int64{},
	}
}

func (t *ProductSalesTable) Set(date string, productId int64, value int) {
	var dateIndex int
	if index, ok := t.dateMap[date]; ok {
		dateIndex = index
	} else {
		// TODO save dateIndex
		dateIndex = len(t.DateLabel)
		t.DateLabel = append(t.DateLabel, date)
		t.dateMap[date] = dateIndex
	}

	var productIndex int
	if index, ok := t.productMap[productId]; ok {
		productIndex = index
	} else {
		productIndex = len(t.ProductLabel)
		t.ProductLabel = append(t.ProductLabel, productId)
		t.productMap[productId] = productIndex
	}

	key := fmt.Sprintf("%d_%d", dateIndex, productIndex)
	t.data[key] = value
}

func (t *ProductSalesTable) Get(date string, productId int64) int {
	dateIndex, ok := t.dateMap[date]
	if !ok {
		return 0
	}
	productIndex, ok := t.productMap[productId]
	if !ok {
		return 0
	}
	key := fmt.Sprintf("%d_%d", dateIndex, productIndex)
	return t.data[key]
}

func (t *ProductSalesTable) ProductMap() map[int64]bool {
	labels := make(map[int64]bool, len(t.productMap))
	for k, _ := range t.productMap {
		labels[k] = true
	}
	return labels
}

func (t *ProductSalesTable) SortedDateLabel() []string {
	keys := make([]string, len(t.DateLabel))
	copy(keys, t.DateLabel)
	fmt.Println("8888888888888888888888888888")
	fmt.Println(keys)
	sort.Strings(keys)
	return keys
}

func (t *ProductSalesTable) SumDate(date string) int {
	dateIndex, ok := t.dateMap[date]
	if !ok {
		return 0
	}
	var sumDate int
	for productIndex, _ := range t.ProductLabel {
		key := fmt.Sprintf("%d_%d", dateIndex, productIndex)
		if value, ok := t.data[key]; ok {
			if value > 0 {
				sumDate += value
			}
		}
	}
	return sumDate
}

func (t *ProductSalesTable) SumProduct(productId int64) int {
	productIndex, ok := t.productMap[productId]
	if !ok {
		return 0
	}
	var sumProduct int
	for dateIndex, _ := range t.DateLabel {
		key := fmt.Sprintf("%d_%d", dateIndex, productIndex)
		if value, ok := t.data[key]; ok {
			if value > 0 {
				sumProduct += value
			}
		}
	}
	return sumProduct
}

func NewTestProductSalesTable() *ProductSalesTable {
	newt := NewProductSalesTable()
	newt.Set("2015-01-01", 56, 9)
	newt.Set("2015-01-01", 55, 9)
	newt.Set("2015-01-01", 52, 9)
	newt.Set("2015-01-01", 51, 9)
	newt.Set("2015-02-03", 50, 9)
	newt.Set("2015-02-03", 54, 9)
	newt.Set("2015-02-04", 55, 9)
	return newt
}

// type SettleAccountTable struct {
// 	data         [][]int
// 	dateMap      map[string]int
// 	productMap   map[int64]int
// 	dateLabel    []string
// 	productLabel []string
// }

// func NewSettleAccountTable() *SettleAccountTable {
// 	return &SettleAccountTable{
// 		data:         [][]int{},
// 		dateMap:      map[string]int{},
// 		productMap:   map[int64]int{},
// 		dateLabel:    []string{},
// 		productLabel: []string{},
// 	}
// }

// func NewTestSettleAccountTable(start, end time.Time) *SettleAccountTable {
// 	r := &SettleAccountTable{
// 		data: [][]int{
// 			{33, 0, 0, 25},
// 			{3333, 3, 0, 25},
// 		},
// 	}
// 	r.dateMap = map[string]int{
// 		"2016-01-01": 0,
// 		"2016-01-02": 1,
// 		"2016-01-03": 2,
// 	}

// 	return r
// }

// func (t *SettleAccountTable) Set(date string, productId int64, value int) {
// 	var dateIndex int
// 	if index, ok := t.dateMap[date]; ok {
// 		dateIndex = index
// 	} else {
// 		// TODO save dateIndex
// 		dateIndex = len(t.dateLabel)
// 		t.dateLabel = append(t.dateLabel, date)
// 		t.dateMap[date] = dateIndex
// 	}

// 	var productId int
// }
