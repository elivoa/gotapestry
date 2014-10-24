package orderservice

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/elivoa/got/debug"
	"github.com/elivoa/gxl"
	"math"
	"strconv"
	"syd/dal"
	"syd/dal/accountdao"
	"syd/dal/orderdao"
	"syd/dal/persondao"
	"syd/model"
	"syd/service/productservice"
	"time"
)

// --------------------------------------------------------------------------------

type OrderDetailJson struct {
	Orders   []int                              `json:"order"`
	Products map[string]*ProductDetalJsonStruct `json:"products"`
}

// todo rename
type ProductDetalJsonStruct struct {
	Id           int             `json:"id"` // product id
	ProductId    string          `json:"pid"`
	Name         string          `json:"name"`
	SellingPrice float64         `json:"price"`
	ProductPrice float64         `json:"productPrice"`
	Colors       []string        `json:"colors"`
	Sizes        []string        `json:"sizes"`
	Quantity     [][]interface{} `json:"quantity"`
	Note         string          `json:"note"`
}

// --------------------------------------------------------------------------------
// TODO: 如此多的方法，还是弄一个类似于Params的东西来接收可变参数。
func CountOrder(status string) (int, error) {
	return orderdao.CountOrder(status)
}

func ListOrder(status string) ([]*model.Order, error) {
	return orderdao.ListOrder(status)
}

func ListOrderPager(status string, limit int, n int) ([]*model.Order, error) {
	return orderdao.ListOrderPager(status, limit, n)
}

func ListOrderByType(orderType model.OrderType, status string) ([]*model.Order, error) {
	return orderdao.ListOrderByType(orderType, status)
}

// >> copied to service.Order
func CreateOrder(order *model.Order) error {
	_processOrderCustomerPrice(order)
	_calculateOrder(order)
	return orderdao.CreateOrder(order)
}

// >> copied to service.Order
func UpdateOrder(order *model.Order) (int64, error) {
	_processOrderCustomerPrice(order)
	_calculateOrder(order)
	return orderdao.UpdateOrder(order)
}

// >> copied to service.Order
// calculate order by type
func _calculateOrder(order *model.Order) {
	switch model.OrderType(order.Type) {
	case model.SubOrder, model.Wholesale:
		if order.Details != nil && len(order.Details) > 0 {
			order.CalculateOrder()
		}
	case model.ShippingInstead:
		// this type of order's total price is calculated by sub
		// orders, which is difficult to calculate, so I calclate sum
		// in page, and then submit to the parent order. So, here do
		// nothing.
	}
}

// >> copied to service.Order
// save customerized price for order
func _processOrderCustomerPrice(order *model.Order) {
	if order.Details == nil {
		return
	}
	sets := map[int]bool{}
	for _, detail := range order.Details {
		if _, ok := sets[detail.ProductId]; ok {
			continue
		}
		if detail.ProductId == 0 { // pass invalid detail item
			continue
		}
		sets[detail.ProductId] = true

		product := productservice.GetProduct(detail.ProductId)
		if product == nil {
			panic(fmt.Sprint("Can not find product ", detail.ProductId))
		}
		if detail.SellingPrice != product.Price {
			// if different, update
			cp := dal.GetCustomerPrice(order.CustomerId, detail.ProductId)
			if cp == nil || cp.Price != detail.SellingPrice {
				if err := dal.SetCustomerPrice(order.CustomerId, detail.ProductId,
					detail.SellingPrice); err != nil {
					panic(err.Error())
				}
			}
		}
	}
}

func CancelOrder(trackNumber int64) error {
	return ChangeOrderStatus(trackNumber, "canceled")
}

func ChangeOrderStatus(trackNumber int64, status string) error {
	rowsAffacted, err := orderdao.UpdateOrderStatus(trackNumber, status)
	if err != nil {
		return err
	}
	if rowsAffacted == 0 {
		return errors.New("No rows affacted!")
	}
	return nil
}

// >> copied to service.order
func GetOrder(id int) (*model.Order, error) {
	return orderdao.GetOrder("id", id)
}

func GetOrderByTrackingNumber(trackingNumber int64) (*model.Order, error) {
	return orderdao.GetOrder("track_number", trackingNumber)
}

func DeleteOrder(trackNumber int64) (affacted int64, err error) {
	affacted, err = orderdao.DeleteOrder(trackNumber)
	return
}

// load all suborders, with all details. set to order
func LoadSubOrders(order *model.Order) ([]*model.Order, error) {
	// now := time.Now()
	suborders, err := orderdao.ListSubOrders(order.TrackNumber)
	if err != nil {
		return nil, err
	}
	// load all details. cascaded.
	for _, o := range suborders {
		details, err := orderdao.GetOrderDetails(o.TrackNumber)
		if err != nil {
			return nil, err
		}
		o.Details = details
	}
	// fmt.Println()
	return suborders, nil
}

// 批量结款
func BatchCloseOrder(money float64, customerId int) {
	debug.Log("Incoming Money: %v", money)
	person, err := persondao.Get(customerId)
	if err != nil {
		panic(err.Error())
	}
	// get unclosed orders for somebody
	orders, err := orderdao.DeliveringUnclosedOrdersByCustomer(customerId)
	if err != nil {
		panic(err.Error())
	}

	// collect totalorder price
	var totalOrderPrice float64
	for _, o := range orders {
		totalOrderPrice += o.SumOrderPrice()
	}

	// money used as total shouldbe: inputmoney + (accountballance - allorder's price)
	totalmoney := money + (person.AccountBallance + totalOrderPrice)

	for _, order := range orders {
		if totalmoney-order.SumOrderPrice() >= 0 {
			err := ChangeOrderStatus(order.TrackNumber, "done")
			if err != nil {
				panic(err.Error())
			}
			totalmoney -= order.SumOrderPrice()
		}
	}
	accountdao.CreateIncoming(&model.AccountIncoming{
		CustomeId: person.Id,
		Incoming:  money,
	})
	// modify customer's accountballance
	person.AccountBallance += money
	persondao.Update(person)

	// create chagne log at the same time:
	accountdao.CreateAccountChangeLog(&model.AccountChangeLog{
		CustomerId: person.Id,
		Delta:      money,
		Account:    person.AccountBallance,
		Type:       2, // create order
		// RelatedOrderTN: 0,
		Reason: "Batch insert",
	})

}

// ________________________________________________________________________________
// ProductJson generator
func OrderDetailsJson(order *model.Order) *OrderDetailJson {
	orders := []int{}
	products := map[string]*ProductDetalJsonStruct{}

	if order.Details != nil {
		for _, detail := range order.Details {
			if detail == nil || detail.ProductId == 0 {
				continue
			}

			// add to cache
			jsonStruct, ok := products[strconv.Itoa(detail.ProductId)]
			if !ok {
				// get product
				product := productservice.GetProduct(detail.ProductId)
				// product, err := productdao.Get(detail.ProductId)
				if product == nil {
					panic("can not find product")
				}

				jsonStruct = &ProductDetalJsonStruct{
					Id:           product.Id,
					ProductId:    product.ProductId,
					Name:         product.Name,
					SellingPrice: detail.SellingPrice,
					ProductPrice: product.Price,
					Colors:       product.Colors,
					Sizes:        product.Sizes,
					Quantity:     [][]interface{}{},
					Note:         detail.Note,
				}
				products[strconv.Itoa(detail.ProductId)] = jsonStruct
				orders = append(orders, product.Id)
			}

			// update quantity
			jsonStruct.Quantity = append(jsonStruct.Quantity,
				[]interface{}{detail.Color, detail.Size, detail.Quantity})
		}
	}
	r := OrderDetailJson{Orders: orders, Products: products}
	// bytes, err := json.Marshal(r)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// return string(bytes)
	return &r
}

func ListOrderByTime(start, end time.Time) ([]*model.Order, error) {
	return orderdao.ListOrderByTime(start, end)
}

// CombineOrderDetails combines OrderDetails into one order, others are ignored.
func CombineOrderDetials(orders ...*model.Order) *model.Order {
	finalOrder := &model.Order{}
	if nil == orders || len(orders) == 0 {
		return finalOrder
	}
	var (
		totalQuantity int
		totalPrice    float64
		odmap         = make(map[string]*model.OrderDetail)
	)

	// can't combined:
	//   detail.SellingPrice

	for idx, o := range orders {
		if o.Details == nil || len(o.Details) == 0 {
			continue
		}
		for _, d := range o.Details {
			odkey := fmt.Sprintf("%v__%v__%v", d.ProductId, d.Color, d.Size)
			detail, ok := odmap[odkey]
			if !ok {
				detail = &model.OrderDetail{
					ProductId:    d.ProductId,
					Color:        d.Color,
					Size:         d.Size,
					Quantity:     0,
					SellingPrice: d.SellingPrice,
					Unit:         d.Unit,
				}
				odmap[odkey] = detail
			}
			detail.Quantity += d.Quantity
			// sum
			totalQuantity += d.Quantity
			totalPrice += float64(d.Quantity) * d.SellingPrice
		}

		// order things
		finalOrder.Note += fmt.Sprint(o.TrackNumber, "; ")
		if finalOrder.ExpressFee > 0 {
			finalOrder.ExpressFee += o.ExpressFee
		}
		finalOrder.DeliveryMethod = o.DeliveryMethod
		if o.DeliveryTrackingNumber != "" {
			finalOrder.DeliveryTrackingNumber += o.DeliveryTrackingNumber + "; "
		} else {
			finalOrder.DeliveryTrackingNumber += "【单号欠缺】; "
		}
		// Accumulated we choose the bigest, instead of sum them.
		if idx == 0 {
			finalOrder.Accumulated = o.Accumulated
		} else {
			finalOrder.Accumulated = math.Min(finalOrder.Accumulated, o.Accumulated)
		}
	}

	// set to final order.
	finalOrder.TotalCount = totalQuantity
	finalOrder.TotalPrice = totalPrice

	finalOrder.Details = make([]*model.OrderDetail, len(odmap))
	for _, d := range odmap {
		finalOrder.Details = append(finalOrder.Details, d)
	}
	return finalOrder
}

func LoadDetails(orders []*model.Order) error {
	if orders != nil {
		for _, o := range orders {
			details, err := orderdao.GetOrderDetails(o.TrackNumber)
			if err != nil {
				return err
			}
			o.Details = details
		}
	}
	return nil
}

func GenerateLeavingMessage(customerId int, date time.Time) (*model.Order, string) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 0, 1)
	orders, err := orderdao.ListOrderByCustomer_Time(customerId, start, end)
	if err != nil {
		panic(err.Error())
		// return err.Error()
	}
	return CombinedLeavingMessage(orders...)
}

func CombinedLeavingMessage(orders ...*model.Order) (*model.Order, string) {
	if orders == nil || len(orders) == 0 {
		return nil, "<<今日无订单!>>"
	}
	neworders := []*model.Order{}
	for _, o := range orders {
		if o != nil {
			if o.Type == uint(model.Wholesale) {
				neworders = append(neworders, o)
			}
		}
	}
	LoadDetails(neworders)
	bigOrder := CombineOrderDetials(neworders...)
	return bigOrder, LeavingMessage(bigOrder)
}

func LeavingMessage(bigOrder *model.Order) string {
	var msg bytes.Buffer
	jo := OrderDetailsJson(bigOrder)
	var sumTotal float64
	var sumQuantity int
	for _, id := range jo.Orders {
		productJson := jo.Products[strconv.Itoa(id)]
		// 例如：奢华宝石
		msg.WriteString(productJson.Name)
		totalQuantity := 0
		for _, q := range productJson.Quantity {
			totalQuantity += q[2].(int)
			sumTotal += float64(q[2].(int)) * productJson.SellingPrice
		}
		sumQuantity += totalQuantity
		// eg: 1件
		msg.WriteString(strconv.Itoa(totalQuantity))
		msg.WriteString("件")

		// details
		if len(productJson.Quantity) >= 1 {
			msg.WriteString("(")
			i := 0
			for _, q := range productJson.Quantity {
				if i > 0 {
					msg.WriteString(", ")
				}
				i += 1
				_color := q[0].(string)
				_size := q[1].(string)
				if _color != "默认颜色" {
					msg.WriteString(_color)
				}
				if _size != "均码" {
					msg.WriteString(_size)
				}
				msg.WriteString(" ")
				msg.WriteString(strconv.Itoa(q[2].(int)))
			}
			msg.WriteString(")")
		}
		msg.WriteString("，")

		// price eg: xxx元
		msg.WriteString(fmt.Sprint(productJson.SellingPrice * float64(totalQuantity)))
		// msg.WriteString(gxl.FormatCurrency(productJson.SellingPrice*float64(totalQuantity), 2))
		msg.WriteString("元")
		msg.WriteString("；")
	}

	// 共计 n件 x元
	msg.WriteString("共计")
	msg.WriteString(strconv.Itoa(sumQuantity))
	msg.WriteString("件")
	msg.WriteString(gxl.FormatCurrency(sumTotal, 2))
	msg.WriteString("元")
	msg.WriteString("；")

	// shipping. sum multi
	switch bigOrder.DeliveryMethod {
	case "SF":
		msg.WriteString("顺风运费")
	case "YTO":
		msg.WriteString("圆通运费")
	case "Depoon":
		msg.WriteString("德邦物流")
	case "Freight":
		msg.WriteString("货运物流")
	default:
		msg.WriteString("【" + bigOrder.DeliveryMethod + "】 运费")
	}

	if bigOrder.ExpressFee > 0 {
		msg.WriteString(fmt.Sprint(bigOrder.ExpressFee))
		msg.WriteString("元；")
	}
	msg.WriteString("单号")
	msg.WriteString(bigOrder.DeliveryTrackingNumber)
	// msg.WriteString("; ")

	// 总计
	msg.WriteString("总计")
	if bigOrder.ExpressFee > 0 {
		msg.WriteString(gxl.FormatCurrency(sumTotal+float64(bigOrder.ExpressFee), 2))
	} else {
		msg.WriteString(gxl.FormatCurrency(sumTotal, 2))
	}
	msg.WriteString("元")
	msg.WriteString("；")

	// 累计欠款
	if bigOrder.Accumulated > 0 {
		// TODO 多个订单的时候累积欠款是错的。
		msg.WriteString("累计欠款：")
		msg.WriteString(fmt.Sprint(bigOrder.Accumulated))
		msg.WriteString(" + ")
		msg.WriteString(fmt.Sprint(int64(sumTotal) + bigOrder.ExpressFee))
		msg.WriteString(" = ")
		msg.WriteString(gxl.FormatCurrency(float64(int64(sumTotal)+bigOrder.ExpressFee)+bigOrder.Accumulated, 2))
		msg.WriteString("元")
		msg.WriteString("；")
	}
	return msg.String()
}
