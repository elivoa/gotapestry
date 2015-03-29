package model

//
// stat by day.
//

type SumStat struct {
	Id         int
	NOrder     int
	NSold      int
	AvgPrice   float64
	TotalPrice float64
}

//
// hot sales model
//

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
