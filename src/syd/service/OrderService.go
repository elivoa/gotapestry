package service

import (
	"errors"
	"fmt"
	"syd/dal"
	"syd/dal/accountdao"
	"syd/dal/orderdao"
	"syd/model"
	"syd/service/personservice"

	"github.com/elivoa/got/db"
	"github.com/elivoa/got/debug"
)

var debug_always_update_order_accumulated = true
var enableSales = true

type OrderService struct{}

func (s *OrderService) EntityManager() *db.Entity {
	return orderdao.EntityManager()
}

func (s *OrderService) GetOrder(id int) (*model.Order, error) {
	return orderdao.GetOrder("id", id)
}

func (s *OrderService) GetOrderByTrackingNumber(trackNumber int64) (*model.Order, error) {
	return orderdao.GetOrder("track_number", trackNumber)
}

// Update Order 的所有逻辑都在这里了；
func (s *OrderService) UpdateOrder(order *model.Order) (*model.Order, error) {

	fmt.Println("0000", order.Details)
	for idx, d := range order.Details {
		fmt.Println("\t00", idx, d, d.SellingPrice)
		fmt.Println("\t00", idx, d.SellingPrice)
		fmt.Println("\t00", idx, d.DiscountPercent)
	}

	// 更新 Tracking Number 到子订单的 Tracking No
	for _, detail := range order.Details {
		// TODO: additional check if order tracking number not match;
		detail.OrderTrackNumber = order.TrackNumber
	}

	var needUpdateBallance = false
	// If status change from other to takeaway, mark as need update ballance.
	// Order中只包含新数据，旧的数据必须从数据库中取出了；也顺便再次验证是否存在这个订单；
	if oldOrder, err := s.GetOrder(order.Id); err != nil {
		return nil, err // panic(err)
	} else {
		// this makes mistake.
		// if oldOrder.DeliveryMethod == "TakeAway" {
		// 	return nil, errors.New("自提状态的订单不能再修改！")
		// }

		// 自提，非自提；状态是发货中，和其他
		// 更新情况1： 修改订单时，由非自提状态改为自提状态，需要修改累计欠款;
		if oldOrder.DeliveryMethod != "TakeAway" && order.DeliveryMethod == "TakeAway" {
			order.Status = "delivering"
			needUpdateBallance = true
		}
		// 更新情况2： 非自提状态点发货，}
		if order.Status == "delivering" { // 需要修改累计欠款
			if oldOrder.Status == "done" || oldOrder.Status == "canceled" {
				return nil, errors.New(fmt.Sprintf("%s状态的订单状态不能修改为delivering！", oldOrder.Status))
			}
			needUpdateBallance = true
		}

		// update order
		_processOrderCustomerPrice(order)    // update order custumize price when confirm order.
		_calculateOrder(order)               // calculate order statistic fields.
		_processingUpdateOrderDetails(order) // update order detail into db;

		_update_order_accumulated(order) // debug option
		// update order
		if _, err := orderdao.UpdateOrder(order); err != nil {
			return nil, err
		}
	}

	// update account ballance. upate stocks left.
	if needUpdateBallance {
		// 代发的父订单不参与非累计欠款, 代发订单由其子订单负责参与累计欠款的统计；
		// 代发父订单的逻辑应该也到不了这里。
		fmt.Println(">> update account ballance. delta: ", -order.SumOrderPrice())

		switch model.OrderType(order.Type) {
		case model.Wholesale, model.SubOrder:

			Account.UpdateAccountBalance(order.CustomerId, -order.SumOrderPrice(),
				"Create Order", order.TrackNumber)

		}

		// update stocks
		_reduceProductStocks(order.Details)

	}
	return order, nil
}

func _update_order_accumulated(order *model.Order) error {
	if debug_always_update_order_accumulated {

		// always update order.accumulated
		// 4. get person, check if customer exists.
		customer, err := Person.GetPersonById(order.CustomerId)
		if err != nil {
			return err
		} else if customer == nil {
			return errors.New(fmt.Sprintf("Customer not found for order! id %v", order.CustomerId))
		}

		// 5. 设置累计欠款；
		order.Accumulated = -customer.AccountBallance
	}
	return nil
}

func _processingUpdateOrderDetails(order *model.Order) error {
	if nil == order {
		return nil // || order.Details == nil || len(order.Details) <= 0
	}
	var createGroup = []*model.OrderDetail{}
	var updateGroup = []*model.OrderDetail{}
	var deleteGroup = []*model.OrderDetail{}

	// 1. load all details;
	details, err := orderdao.GetOrderDetails(order.TrackNumber)
	if err != nil {
		return err
	}
	if details == nil || len(details) == 0 {
		if order.Details != nil {
			createGroup = order.Details
		}
	} else {
		// normal case, create, update and delete;
		var deleteWhoIsFalse = make([]bool, len(details)) //  make(map[int]bool, len(details))
		if order.Details != nil {
			// find who should create and who need update.
			for _, d := range order.Details {
				// Liner find match in details.
				var find = false
				for idx2, d2 := range details {
					// find the matched item.
					if d.ProductId == d2.ProductId && d.Color == d2.Color && d.Size == d2.Size {
						// assign id into it; update operation need id.
						d.Id = d2.Id
						// if any value changed; If quantity changed to 0, delete it;
						if d.Quantity != d2.Quantity || d.SellingPrice != d2.SellingPrice || d.Note != d2.Note {
							if d.Quantity > 0 {
								updateGroup = append(updateGroup, d)
							}
						}
						// if quantity equails 0, mark as delete;
						find = true
						if d.Quantity == 0 {
							deleteWhoIsFalse[idx2] = false // nothing, just remind it.
						} else {
							deleteWhoIsFalse[idx2] = true
						}
						break
					}
				}
				if !find {
					createGroup = append(createGroup, d) // if not found, create this.
				}
			}
		}

		// --------------------------------------------------------------------------------
		fmt.Println(">>>> details and order details:")
		if nil != details {
			for _, d := range details {
				fmt.Println("\tdetails: ", d.OrderTrackNumber, d.Color, d.Size, " = ", d.Quantity, d.SellingPrice)
			}
		}
		if nil != order.Details {
			for _, d := range order.Details {
				fmt.Println("\torder details: ", d.OrderTrackNumber, d.Color, d.Size, " = ", d.Quantity, d.SellingPrice)
			}
		}
		fmt.Println(">>>> who is false?")
		for idx, b := range deleteWhoIsFalse {
			fmt.Println("\t >> who is false: ", idx, b)
		}

		// who will be deleted?
		for idx, b := range deleteWhoIsFalse { // } i := 0; i < len(details); i++ {
			if !b {
				deleteGroup = append(deleteGroup, details[idx])
			}
		}
	}

	var debugdetails = false
	if debugdetails {
		fmt.Println("\n\n\n--------------------------------------------------------------------------------")
		fmt.Println("Order Detail Create Group:")
		if nil != createGroup {
			for _, d := range createGroup {
				fmt.Println("\tcreate: ", d.OrderTrackNumber, d.Color, d.Size, " = ", d.Quantity, d.SellingPrice)
			}
		}
		fmt.Println("Order Detail Update Group:")
		if nil != updateGroup {
			for _, d := range updateGroup {
				fmt.Println("\tupdate: ", d.OrderTrackNumber, d.Color, d.Size, " = ", d.Quantity, d.SellingPrice)
			}
		}
		fmt.Println("Order Detail Delete Group:")
		if nil != deleteGroup {
			for _, d := range deleteGroup {
				fmt.Println("\tdelete: ", d.OrderTrackNumber, d.Color, d.Size, " = ", d.Quantity, d.SellingPrice)
			}
		}
		fmt.Println("==========================================================================\n\nDebug Done!")
	}
	// final process: create, update, and delete
	if createGroup != nil {
		if err := orderdao.CreateOrderDetail(createGroup); err != nil {
			return err
		}
	}
	if updateGroup != nil {
		if err := orderdao.BatchUpdateOrderDetail(updateGroup); err != nil {
			return err
		}
	}
	if deleteGroup != nil {
		if err := orderdao.DeleteOrderDetails(deleteGroup); err != nil {
			return err
		}
	}

	return nil
}

// 创建订单的所有逻辑都在这里
func (s *OrderService) CreateOrder(order *model.Order) (*model.Order, error) {
	fmt.Println("0000", order)
	// 这个步骤很重要，判断是否订单已经存在了；如果存在了，需要换一个订单号再试；
	if newtn := makeSureOrderTNUnique(order.TrackNumber); newtn > 0 {
		order.TrackNumber = newtn
	}

	for _, detail := range order.Details {
		detail.OrderTrackNumber = order.TrackNumber
	}

	var needUpdateBallance = false // 需要写入Person的累计欠款

	// If order delivery method is `takeaway`, chagne order status to `delivering` and
	// update account ballance; In other situation change status to `toprint`.
	if order.DeliveryMethod == "TakeAway" {
		order.Status = "delivering"
		needUpdateBallance = true
	} else {
		order.Status = "toprint"
	}

	// organize some data
	_processOrderCustomerPrice(order)
	_calculateOrder(order)

	// create order detail into db;
	if order.Details != nil && len(order.Details) > 0 {
		if err := orderdao.CreateOrderDetail(order.Details); err != nil {
			return nil, err
		}
	}

	// create order into db
	// fmt.Println("======================================")
	// fmt.Println("order.CreateTimeis ", order.CreateTime)
	// fmt.Println("======================================")
	if err := orderdao.CreateOrder(order); err != nil {
		return nil, err
	}

	_update_order_accumulated(order) // debug option

	// update account ballance
	if needUpdateBallance {
		Account.UpdateAccountBalance(order.CustomerId, -order.SumOrderPrice(),
			"Create Order", order.TrackNumber)

		// last step: update stocks. 修改累计欠款的时候就修改库存；
		_reduceProductStocks(order.Details)
	}

	return order, nil
}

// 订单中剪掉库存数量；
func _reduceProductStocks(orderDetails []*model.OrderDetail) error {

	if nil == orderDetails || len(orderDetails) == 0 {
		return nil
	}
	fmt.Println("\n\n========== reduce product stocks")
	for _, detail := range orderDetails {
		if detail == nil {
			continue
		}
		// Log oldStock? newStock?
		_, _, err := Stock.UpdateStockDelta(
			(int64)(detail.ProductId), detail.Color, detail.Size, -detail.Quantity)
		if err != nil {
			return err
		}
	}
	return nil
}

// 发货
func (s *OrderService) DeliverOrder(trackNumber int64, deliveryTrackingNumber, deliveryMethod string,
	expressFee int64) (*model.Order, error) {

	// 1. get order form db.
	order, err := s.GetOrderByTrackingNumber(trackNumber)
	if err != nil {
		return order, err
	}

	// 2. status check
	if order == nil {
		return nil, errors.New(fmt.Sprintf("Order %d doesn't exist!", trackNumber))
	}
	switch order.Status {
	case "delivering", "done", "canceled":
		return order, errors.New(fmt.Sprintf("状态为%s的订单不能发货！", order.Status))
	}

	// 3. set data back to order.
	order.DeliveryTrackingNumber = deliveryTrackingNumber
	order.DeliveryMethod = deliveryMethod
	order.ExpressFee = expressFee
	order.Status = "delivering"

	// 4. get person, check if customer exists.
	customer, err := Person.GetPersonById(order.CustomerId)
	if err != nil {
		return order, err
	} else if customer == nil {
		return order, errors.New(fmt.Sprintf("Customer not found for order! id %v", order.CustomerId))
	}

	// 5. 设置累计欠款；
	order.Accumulated = -customer.AccountBallance

	// 5. save order changes.
	if _, err := s.UpdateOrder(order); err != nil {
		return order, err
	}
	return order, nil
}

// Make sure order tracking number not conflict, by checking if tn exists,
// change to another but not checked again. So, it's limited.
func makeSureOrderTNUnique(tn int64) int64 {
	if order, err := orderdao.GetOrder("track_number", tn); err != nil {
		panic(err)
	} else if order == nil {
		return 0 // success
	} else {
		return model.GenerateOrderId() // try only once.
	}
}

// calculate order,
func _calculateOrder(order *model.Order) {
	switch model.OrderType(order.Type) {
	case model.SubOrder, model.Wholesale:
		if order.Details != nil && len(order.Details) > 0 {
			order.CalculateOrder()
		}
	case model.ShippingInstead:
		// this type of order's total price is calculated by sub
		// orders, which is difficult to calculate, so I calclate sum
		// in page, and then submit to the parent order.
		// So, here does nothing.
	}
}

var log_processOrderCustomerPrice = true

// save customerized price for order
func _processOrderCustomerPrice(order *model.Order) {
	if order.Details == nil {
		return
	}
	sets := map[int]bool{}
	for _, detail := range order.Details {
		if _, ok := sets[detail.ProductId]; ok { // make sure one detail be processed once.
			continue
		}
		if detail.ProductId <= 0 { // pass invalid detail item
			continue
		}
		sets[detail.ProductId] = true

		var needUpdatePrice = false
		customePrice := dal.GetCustomerPrice(order.CustomerId, detail.ProductId)
		if customePrice != nil && customePrice.Price != detail.SellingPrice {
			// has customer price
			needUpdatePrice = true
		} else {
			// don't has customer price, load product price;

			// TODO: performance issue, batch get product.
			if product, err := Product.GetFullProduct(detail.ProductId); err != nil {
				panic(err)
			} else if nil == product {
				panic(fmt.Sprint("Can not find product ", detail.ProductId))
			} else {
				if product.Price != detail.SellingPrice {
					needUpdatePrice = true
				}
			}
		}

		// update selling price.
		if needUpdatePrice {
			if err := dal.SetCustomerPrice(order.CustomerId, detail.ProductId,
				detail.SellingPrice); err != nil {
				panic(err.Error())
			}
		}
		// // here has bugs;
		// fmt.Println("\n\n\n\n\n >>>> price price price price price price price ")
		// fmt.Println(" >>>> ", detail.SellingPrice, product.Price)
		// if detail.SellingPrice != product.Price {
		// 	// if different, update
		// 	cp := dal.GetCustomerPrice(order.CustomerId, detail.ProductId)
		// 	// fmt.Println(">>>>>>>>>>>> ", cp)
		// 	if cp == nil || cp.Price != detail.SellingPrice {
		// 	}
		// }
	}
}

func (s *OrderService) ListOrders(parser *db.QueryParser, withs Withs) ([]*model.Order, error) {
	if orders, err := orderdao.ListOrders(parser); err != nil {
		return nil, err
	} else {
		// TODO: Print warrning information when has unused withs.
		if withs&WITH_PERSON > 0 {
			if err := s.FillOrderSlicesWithPerson(orders); err != nil {
				return nil, err
			}
		}
		return orders, nil
	}
}

// orderlist is passed by pointer.
func (s *OrderService) FillOrderSlicesWithPerson(orders []*model.Order) error {
	var idset = map[int64]bool{}
	for _, order := range orders {
		idset[int64(order.CustomerId)] = true
	}
	personmap, err := Person.BatchFetchPersonByIdMap(idset)
	if err != nil {
		return err
	}
	if nil != personmap && len(personmap) > 0 {
		for _, order := range orders {
			if person, ok := personmap[int64(order.CustomerId)]; ok {
				order.Customer = person
			}
		}
	}
	return nil
}

// --------------------------------------------------------------------------------
// special

//
func (s *OrderService) BatchCloseOrder(money float64, customerId int) {
	debug.Log("Incoming Money: %v", money)
	person, err := Person.GetPersonById(customerId)
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
			err := s.ChangeOrderStatus(order.TrackNumber, "done")
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
	personservice.Update(person) // TODO: chagne place

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

func (s *OrderService) ChangeOrderStatus(trackNumber int64, status string) error {
	rowsAffacted, err := orderdao.UpdateOrderStatus(trackNumber, status)
	if err != nil {
		return err
	}
	if rowsAffacted == 0 {
		return errors.New("No rows affacted!")
	}
	return nil
}
