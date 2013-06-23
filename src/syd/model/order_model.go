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
	TrackNumber    int64
	Status         string "" // New | Closed
	DeliveryMethod string "" // TakeAway | Express(圆通，顺风)

	CustomerId int // reference

	// summarization
	TotalPrice float64
	TotalCount int
	PriceCut   float64

	// Not in db
	Details []*OrderDetail `inject:"slice"`

	Note string

	// times
	CreateTime time.Time
	UpdateTime time.Time
	CloseTime  time.Time
}

func NewOrder() *Order {
	order := &Order{
		TrackNumber: GenerateOrderId(), DeliveryMethod: "Express",
		CreateTime: time.Now(),
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
  Order Details
*/
type OrderDetail struct {
	Id int
	// reference
	OrderTrackNumber int64 // reference
	ProductId        int   // reference
	Quantity         int
	Unit             string  // always 件, no use
	SellingPrice     float64 //售价
	Note             string
}

// check avaliability
func (d *OrderDetail) IsValid() bool {
	if d.ProductId == 0 && d.Quantity == 0 && d.SellingPrice == 0 && d.Note == "" {
		return false
	} else {
		return true
	}
}
func (d *OrderDetail) String() string {
	return fmt.Sprintf("OrderDetail:Id:%v, OrderTrackNumber:%v, ProductId:%v, Note:%v",
		d.Id, d.OrderTrackNumber, d.ProductId, d.Note)
}



















