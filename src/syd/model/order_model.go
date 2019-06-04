package model

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type OrderType uint

// 不可以改变顺序，数据库中是按照增量存储的。
const (
	Wholesale       OrderType = iota // 0 - 大货
	ShippingInstead                  // 1 - 代发
	SubOrder                         // 2 - 子订单（主要用于 1-代发）// 不应该包含在数量统计中。
)

type OrderStatus uint

// TODO 这个没用了，数据库里面存储的是字符串的；
const (
	ToPrint   OrderStatus = iota // new order
	ToDeliver                    //
	Delivering
	Done
	Canceled
)

/*________________________________________________________________________________
  Order Status DisplayMap
  TODO: a better place to show this?
*/
var OrderStatusDisplayMap = map[string]string{
	"toprint":    "待打印",  // 新订单默认状态， 等待打印
	"todeliver":  "待发货",  // 大货订单打印之后进入代发货状态，代发订单进入正在发货状态
	"delivering": "正在发货", // 打印订单之后，转为发货状态。并且取快照状态的累计欠款
	"done":       "已完成",  // 完成订单。
	"canceled":   "已取消",  // 已取消订单，累计欠款不显示。

	"new": "未处理", // 退货订单[new:未处理, done:已完成（换货或退款）]
}

// ________________________________________________________________________________
// Order 的Status 涉及到累计欠款，因此状态绝对不能乱改。必须严格按照流程走。
//
type Order struct {
	Id          int
	TrackNumber int64  `` // real identification
	Status      string `` // todeliver | delivering | done | canceled | (all)
	Type        uint   `` //OrderType  // 代发 | 大货 | 子订单
	CustomerId  int    // TODO: change this into int64

	// shipping info
	DeliveryMethod         string `` // YTO, SF, Depoon, Freight, TakeAway
	DeliveryTrackingNumber string `` // 快递单号
	ExpressFee             int64  `` // -1 means 到付
	ShippingAddress        string `` // this only used in ShippingInstead

	// price summarization.
	TotalPrice         float64 // not include expressfee
	TotalCount         int
	OriginalTotalPrice float64 // order's original total price, before discount.
	PriceCut           float64 // currently not used.
	Accumulated        float64 // 上期累计欠款快照（不包含本期订单价格以及代发费用）

	Note              string
	ParentTrackNumber int64 `` // if has value it's a suborder

	Details []*OrderDetail `inject:"slice"` // cascated

	// times
	CreateTime time.Time
	UpdateTime time.Time
	CloseTime  time.Time

	// additional containers
	Customer *Person
}

// 根据结构，这个应该设计成两个表吧CS信息独立出去，为了省市，重复价格备注字段。
type OrderDetail struct {
	Id int
	// reference
	OrderTrackNumber int64 // reference
	ProductId        int   // reference
	Color            string
	Size             string
	Quantity         int
	SellingPrice     float64 // 售价，实际售价。
	ProductPrice     float64 // 商品原价.
	DiscountPercent  int     // 当时的折扣百分比.

	Unit string // always 件, Not Used Yet.
	Note string

	// additional containers
	Product *Product
}

func NewOrder() *Order {
	order := &Order{
		TrackNumber: GenerateOrderId(),
		Status:      "toprint",
		CreateTime:  time.Now(),
	}
	order.Details = []*OrderDetail{
		&OrderDetail{},
	}
	return order
}

func main() {
	fmt.Println(33)
	for i := 0; i <= 10; i++ {
		fmt.Println(GenerateOrderId())
	}
}

func GenerateOrderId() int64 {
	now := time.Now()
	// var newid int64
	y, M, d := now.Date()
	h, m, s := now.Clock()
	var first int = 0 +
		(y-2000)*10000000000 +
		int(M)*100000000 +
		d*1000000 +
		h*10000 +
		m*100 +
		s*1 +
		0

	var final int64 = int64(rand.Intn(9999)) + int64(first*10000)

	if f := false; f {
		fmt.Println("first is ", first, s)
	}
	return final
	// if value, err := strconv.ParseInt(time.Now().Format("0601020304"), 10, 64); err != nil {
	// 	panic(err.Error())
	// } else {
	// 	return value*10000 + int64(rand.Intn(9999))
	// }
}

func (order *Order) DisplayStatus() string {
	if display, ok := OrderStatusDisplayMap[order.Status]; ok {
		return display
	}
	return order.Status
}

var enableSales = true

// CalculateOrder calucate order's price. including discount.
func (order *Order) CalculateOrder() {
	// loop to assign values more
	var (
		count       int     // total count
		sum         float64 // total price
		originalSum float64 // totoal price without discount.
	)
	for idx, d := range order.Details {
		if !d.IsValid() {
			order.Details[idx] = nil
			continue
		}

		count += d.Quantity

		// 高级部分，计算价格。
		var price = d.SellingPrice * float64(d.Quantity)
		var originalPrice = price
		if enableSales && d.DiscountPercent > 0 {
			var discountPrice = float64(d.Quantity) * d.ProductPrice * (float64(d.DiscountPercent) / 100)
			if discountPrice < price {
				originalPrice = price
				price = discountPrice
			}
		}

		sum += price
		originalSum += originalPrice

		// assign tracking number
		d.OrderTrackNumber = order.TrackNumber
	}
	order.TotalPrice = sum
	order.TotalCount = count
	if enableSales {
		order.OriginalTotalPrice = originalSum
	}
}

// ----  Show Helper  ----------------------------------------------------------------------------
func (order *Order) IsStatus(status ...string) bool {
	return order.StatusIs(status...)
}

func (order *Order) StatusIs(status ...string) bool {
	for _, s := range status {
		if s == order.Status {
			return true
		}
	}
	return false
}

func (order *Order) IsDaoFu() bool {
	return order.ExpressFee == -1
}

func (order *Order) DeliveryMethodIs(deliveryMethod string) bool {
	return deliveryMethod == order.DeliveryMethod
}

func (order *Order) TypeIs(t uint) bool {
	return order.Type == t
}

func (order *Order) HasAccumulated() bool {
	return order.Accumulated != 0
}

// Total price + express fee
func (order *Order) SumOrderPrice() float64 {
	var sum float64 = order.TotalPrice
	if order.ExpressFee > 0 {
		sum += float64(order.ExpressFee)
	}
	// return sum
	return round(sum, 2)
	// return math.Floor((sum+0.5)*100) / 100
}

func round(v float64, decimals int) float64 {
	var pow float64 = 1
	for i := 0; i < decimals; i++ {
		pow *= 10
	}
	return float64(int((v*pow)+0.5)) / pow
}

// Total price + express fee
func (order *Order) TotalPriceFixed() float64 {
	var sum float64 = order.TotalPrice
	if order.ExpressFee > 0 {
		sum += float64(order.ExpressFee)
	}

	// return sum
	fmt.Println(">>>>>>>>>>>>>>>>>>", sum, round(sum, 2))
	return round(sum, 2)
}

/*________________________________________________________________________________
  Order Details functions
*/
// check avaliability
func (d *OrderDetail) IsValid() bool {
	if d.ProductId == 0 && d.Quantity == 0 && d.SellingPrice == 0 && d.Note == "" {
		return false
	}
	return true
}

// TODO discount
func (d *OrderDetail) TotalPrice() float64 {
	a := float64(d.Quantity) * d.SellingPrice
	return math.Floor(a + 0.5)
}

func (d *OrderDetail) String() string {
	return fmt.Sprintf("OrderDetail(%v), TN:%v, Product:%v_%v_%v quantity:%v, Note:[%v]",
		d.Id, d.OrderTrackNumber, d.ProductId, d.Color, d.Size, d.Quantity, d.Note)
}
