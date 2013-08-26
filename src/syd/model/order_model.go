package model

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type OrderType uint

const (
	Wholesale OrderType = iota
	ShippingInstead
	SubOrder
)

type OrderStatus uint

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
}

// ________________________________________________________________________________
// Order 的Status 涉及到累计欠款，因此状态绝对不能乱改。必须严格按照流程走。
//
type Order struct {
	Id          int
	TrackNumber int64  `` // real identification
	Status      string `` // todeliver | delivering | done | canceled | (all)
	Type        uint   `` //OrderType  // 代发 | 大货 | 子订单
	CustomerId  int

	// shipping info
	DeliveryMethod         string `` // YTO, SF, DEPPON, FREIGHT, TakeAway
	DeliveryTrackingNumber string `` // 快递单号
	ExpressFee             int64  `` // -1 means 到付
	ShippingAddress        string `` // this only used in ShippingInstead

	// price summarization.
	TotalPrice  float64 // not include expressfee
	TotalCount  int
	PriceCut    float64 // currently not used.
	Accumulated float64 // 上期累计欠款快照（不包含本期订单价格以及代发费用）

	Note              string
	ParentTrackNumber int64 `` // if has value it's a suborder

	Details []*OrderDetail `inject:"slice"` // cascated

	// times
	CreateTime time.Time
	UpdateTime time.Time
	CloseTime  time.Time
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
	SellingPrice     float64 //售价
	Unit             string  // always 件, NoUse

	Note string
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

func GenerateOrderId() int64 {
	value, err := strconv.ParseInt(
		fmt.Sprintf("%v%v", time.Now().Format("0601020304"), rand.Intn(999)),
		10, 64)
	if err != nil {
		panic(err.Error())
	}
	return value
}

func (order *Order) DisplayStatus() string {
	display, ok := OrderStatusDisplayMap[order.Status]
	if ok {
		return display
	} else {
		return order.Status
	}
}

func (order *Order) CalculateOrder() {
	// loop to assign values more
	var (
		sum   float64 = 0
		count int     = 0
	)
	for idx, d := range order.Details {
		if !d.IsValid() {
			order.Details[idx] = nil
			continue
		}

		// sum values
		sum += d.SellingPrice * float64(d.Quantity)
		count += d.Quantity

		// assign tracking number
		d.OrderTrackNumber = order.TrackNumber
	}
	order.TotalPrice = sum
	order.TotalCount = count
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
	return order.Accumulated > 0
}

// Total price + express fee
func (order *Order) SumOrderPrice() float64 {
	var sum float64 = order.TotalPrice
	if order.ExpressFee > 0 {
		sum += float64(order.ExpressFee)
	}
	return sum
}

/*________________________________________________________________________________
  Order Details functions
*/
// check avaliability
func (d *OrderDetail) IsValid() bool {
	if d.ProductId == 0 && d.Quantity == 0 && d.SellingPrice == 0 && d.Note == "" {
		return false
	} else {
		return true
	}
}

func (d *OrderDetail) TotalPrice() float64 {
	return float64(d.Quantity) * d.SellingPrice
}

func (d *OrderDetail) String() string {
	return fmt.Sprintf("OrderDetail(%v), TN:%v, Product:%v_%v_%v quantity:%v, Note:[%v]",
		d.Id, d.OrderTrackNumber, d.ProductId, d.Color, d.Size, d.Quantity, d.Note)
}
