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

// ________________________________________________________________________________
// Order 的Status 涉及到累计欠款，因此状态绝对不能乱改。必须严格按照流程走。
//
type Order struct {
	Id                     int
	TrackNumber            int64  `` // real identification
	Status                 string `` // todeliver | delivering | done | canceled | (all)
	DeliveryMethod         string `` // YTO, SF, TakeAway
	DeliveryTrackingNumber string `` // 快递单号
	ExpressFee             int64  `` // -1 means 到付

	CustomerId int            // reference
	Details    []*OrderDetail `inject:"slice"` // cascated

	// summarization, not user input, calculated. Persisted in DB.
	TotalPrice  float64 // not include expressfee
	TotalCount  int
	PriceCut    float64 // currently not used.
	Accumulated float64 // 上期累计欠款快照（不包含本期订单价格以及代发费用）

	Note string

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
		TrackNumber:    GenerateOrderId(),
		Status:         "todeliver",
		DeliveryMethod: "Express",
		CreateTime:     time.Now(),
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
	// loop to assign valuesmroe
	var (
		sum   float64 = 0
		count int     = 0
	)
	fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	for idx, d := range order.Details {
		fmt.Println(d)
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

func (order *Order) IsStatus(status ...string) bool {
	for _, s := range status {
		if s == order.Status {
			return true
		}
	}
	return false
}

func (order *Order) HasAccumulated() bool {
	return order.Accumulated > 0
}

func (order *Order) SumOrderPrice() float64 {
	var sum float64
	sum += order.TotalPrice
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

func (d *OrderDetail) String() string {
	return fmt.Sprintf("OrderDetail(%v), TN:%v, Product:%v_%v_%v quantity:%v, Note:[%v]",
		d.Id, d.OrderTrackNumber, d.ProductId, d.Color, d.Size, d.Quantity, d.Note)
}

/*________________________________________________________________________________
  Order Status DisplayMap
  TODO: a better place to show this?
*/
var OrderStatusDisplayMap = map[string]string{
	"todeliver":  "待发货",  // 新订单默认状态
	"delivering": "正在发货", // 打印订单之后，转为发货状态。并且取快照状态的累计欠款
	"done":       "已完成",  // 完成订单。
	"canceled":   "已取消",  // 已取消订单，累计欠款不显示。
}
