package model

import ()

//
// 发样衣, 将样衣发到大家手中。
//

// 二维表，
// x:product - int64
// y:selected-customer - int64
// value: - int
//    0: 没发
//    -1: 被标记为不可以发.
//    20150908: 发样衣的日期.
// Note: 只取1自然年内数据.
type ProductBroadcast struct {
	data []int64
}

func NewProductBroadcast() *ProductBroadcast {
	return &ProductBroadcast{
		data: []int64{},
	}
}

func (pb *ProductBroadcast) Set(customerId int64, productId int64, value int) {
	// TO Be Continued....
}

// type HotSaleProduct struct {
// 	ProductId   int64
// 	ProductName string
// 	Sales       int
// 	Specs       map[string]int
// }

// type HotSaleProducts []*HotSaleProduct

// func (p HotSaleProducts) Len() int           { return len(p) }
// func (p HotSaleProducts) Less(i, j int) bool { return p[i].Sales > p[j].Sales }
// func (p HotSaleProducts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
