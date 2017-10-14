package model

import (
	"fmt"
	"time"
)

//
// stat by day.
//

type SumStat struct {
	Id         string
	CreateTime time.Time
	NOrder     int
	NSold      int
	AvgPrice   float64
	TotalPrice float64
}

var EmptySumStat = &SumStat{}

//********************************************************************************
// hot sales model
//********************************************************************************

type HotSales struct {
	HSProduct HotSaleProducts
}

type HotSaleProduct struct {
	ProductId   int64
	ProductName string
	Sales       int
	Specs       map[string]int
}

type HotSaleProducts []*HotSaleProduct

func (p HotSaleProducts) Len() int           { return len(p) }
func (p HotSaleProducts) Less(i, j int) bool { return p[i].Sales > p[j].Sales }
func (p HotSaleProducts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

//********************************************************************************
// product daily sales data.
//********************************************************************************

type SalesNode struct {
	Key   string
	Value int
}

type ProductSales []*SalesNode

func (p ProductSales) Len() int           { return len(p) }
func (p ProductSales) Less(i, j int) bool { return p[i].Value > p[j].Value }
func (p ProductSales) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (p ProductSales) Labels() []string {
	var labels = []string{}
	for _, node := range p {
		labels = append(labels, node.Key)
	}
	return labels
}

func (p ProductSales) Datas() []int {
	var labels = []int{}
	for _, node := range p {
		labels = append(labels, node.Value)
	}
	return labels
}

//********************************************************************************
// Hotsales...Name....
//********************************************************************************

type BestBuyerListItem struct {
	CustomerId   int64
	CustomerName string
	Quantity     int
	SalePrice    float64
	// TotalPrice   float64
}

func (m BestBuyerListItem) TotalPrice() float64 {
	return m.SalePrice * float64(m.Quantity)
}

type BestBuyerList []*BestBuyerListItem

//********************************************************************************
// Profit
//********************************************************************************

type Profits struct {
	Profit ProfitM
}

type ProfitModel struct {
	ProductId    int64
	ProductName  string
	Sales        int // 销量
	Specs        map[string]int
	CustomerId   int64
	CustomerName string
	SellingPrice float64
	FactoryPrice float64
}

var juntan_yunfei float64 = 1

// 这条记录的利润，这条记录的利润率。
func (p ProfitModel) Profit() float64 {
	if p.FactoryPrice == 0 {
		return 0
	} else {
		return (p.SellingPrice - p.FactoryPrice) * float64(p.Sales)
	}
}

func (p ProfitModel) ProfitRate() float64 {
	if p.FactoryPrice > 0 {
		return (p.SellingPrice - p.FactoryPrice - juntan_yunfei) / p.FactoryPrice
	}
	return 0
}

func (p ProfitModel) ProfitRateString() string {
	return fmt.Sprintf("%.1f %%", (p.SellingPrice-p.FactoryPrice-juntan_yunfei)/p.FactoryPrice*100)
}

type ProfitM []*ProfitModel

func (p ProfitM) Len() int           { return len(p) }
func (p ProfitM) Less(i, j int) bool { return p[i].Sales > p[j].Sales }
func (p ProfitM) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// person profit
type PersonProfit struct {
	CustomerId   int64
	CustomerName string
	Profit       float64
	ProfitRate   float64
}

type PersonProfits []*PersonProfit

func (p PersonProfits) Len() int           { return len(p) }
func (p PersonProfits) Less(i, j int) bool { return p[i].ProfitRate > p[j].ProfitRate }
func (p PersonProfits) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
