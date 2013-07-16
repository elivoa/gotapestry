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

type Order struct {
	Id             int
	TrackNumber    int64  `` // real identification
	Status         string `` // todeliver | delivering | done | canceled | (all)
	DeliveryMethod string `` // YTO, SF, TakeAway
	ExpressFee     int64  `` // -1 means 到付

	CustomerId int            // reference
	Details    []*OrderDetail `inject:"slice"` // cascated

	// summarization, not user input, calculated. Persisted in DB.
	TotalPrice float64
	TotalCount int
	PriceCut   float64

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
