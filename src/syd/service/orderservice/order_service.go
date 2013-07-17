package orderservice

import (
	"strconv"
	"syd/dal"
	"syd/dal/orderdao"
	"syd/model"
	"syd/service/productservice"
)

func ListOrder(status string) ([]*model.Order, error) {
	return orderdao.ListOrder(status)
}

func CreateOrder(order *model.Order) error {
	_processOrderCustomerPrice(order)
	return orderdao.CreateOrder(order)
}

func UpdateOrder(order *model.Order) (int64, error) {
	_processOrderCustomerPrice(order)
	return orderdao.UpdateOrder(order)
}

func _processOrderCustomerPrice(order *model.Order) {
	if order.Details == nil {
		return
	}
	sets := map[int]bool{}
	for _, detail := range order.Details {
		if _, ok := sets[detail.ProductId]; ok {
			continue
		}
		sets[detail.ProductId] = true

		product := productservice.GetProduct(detail.ProductId)
		if product == nil {
			panic("can not find product")
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

func GetOrder(id int) (*model.Order, error) {
	return orderdao.GetOrder(id)
}

// ________________________________________________________________________________
// ProductJson generator
func ProductDetailJson(order *model.Order) *OrderDetailJson {
	orders := []int{}
	products := map[string]*ProductDetalJsonStruct{}

	if order.Details != nil {
		for _, detail := range order.Details {
			if detail.ProductId == 0 {
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

type OrderDetailJson struct {
	Orders   []int                              `json:"order"`
	Products map[string]*ProductDetalJsonStruct `json:"products"`
}

type ProductDetalJsonStruct struct {
	Id           int             `json:"id"` // product id
	Name         string          `json:"name"`
	SellingPrice float64         `json:"price"`
	ProductPrice float64         `json:"productPrice"`
	Colors       []string        `json:"colors"`
	Sizes        []string        `json:"sizes"`
	Quantity     [][]interface{} `json:"quantity"`
	Note         string          `json:"note"`
}
