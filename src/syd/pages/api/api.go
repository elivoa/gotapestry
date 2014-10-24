/*
SYD Project

API: webservices open to others.
*/
package api

import (
	"encoding/json"
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/debug"
	"syd/dal"
	"syd/dal/productdao"
	"syd/service/personservice"
	"syd/service/productservice"
)

type Api struct {
	core.Page

	APIName string `path-param:"1"`
	Param1  int    `path-param:"2"`
	Param2  int    `path-param:"3"`
}

func (p *Api) Setup() (string, string) {
	switch p.APIName {
	case "person":
		person := personservice.GetPerson(p.Param1)
		return toJson(person)

	case "product":
		product := productservice.GetProduct(p.Param1)
		return toJson(product)

	case "customer_price":
		return "json", getCustomerPrice(p.Param1, p.Param2)
	}

	return needName()
}

// Helpers
// --------------------------------------------------------------------------------
func getCustomerPrice(personId int, productId int) string {
	// TODO extract to service
	var price float64 = -1        // customer price
	var productPrice float64 = -1 // product price
	if personId > 0 {
		// get customer price
		customerPrice := dal.GetCustomerPrice(personId, productId)
		fmt.Println("\n\n\n\n\n >>>> price price price price price price price ")
		fmt.Println(" >>>> ", customerPrice)
		if nil != customerPrice {
			price = customerPrice.Price
		}
	}

	// get product price
	product, err := productdao.Get(productId)
	if err == nil && product != nil {
		productPrice = product.Price
	}
	if price <= 0 {
		price = productPrice
	}
	return fmt.Sprintf("{\"price\":%v, \"productPrice\":%v}", price, productPrice)
}

// Helper error return functions.
// --------------------------------------------------------------------------------
func toJson(obj interface{}) (string, string) {
	if obj == nil {
		return notFound()
	} else {
		personbytes, err := json.Marshal(obj)
		if err != nil {
			return onError(err)
		}
		return "json", string(personbytes)
	}
}

func notFound() (string, string) {
	return "json", "{error: 'No element Found'}"
}

func needName() (string, string) {
	return "json", "{error: 'API Function name must be specified.'}"
}

func onError(err error) (string, string) {
	debug.Error(err)
	return "json", fmt.Sprintf("{error: '%v'}", err.Error())
}
