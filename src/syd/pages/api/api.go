/*
SYD Project

API: webservices open to others.
*/
package api

import (
	"encoding/json"
	"fmt"
	"got/core"
	"got/debug"
	"got/register"
	"syd/dal"
	"syd/service/productservice"
)

func Register() {
	register.Page(Register,
		&Api{},
	)
}

type Api struct {
	core.Page

	APIName string `path-param:"1"`
	Param1  int    `path-param:"2"`
	Param2  int    `path-param:"3"`
}

func (p *Api) Setup() (string, string) {
	switch p.APIName {
	case "person":
		person := dal.GetPerson(p.Param1)
		return toJson(person)

	case "product":
		product := productservice.GetProduct(p.Param1)
		fmt.Println("********************************************************************************")
		fmt.Println(product)
		fmt.Println(product.Colors)
		fmt.Println(product.Sizes)
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
		if nil != customerPrice {
			price = customerPrice.Price
		}
	}

	// get product price
	product, err := dal.GetProduct(productId)
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
