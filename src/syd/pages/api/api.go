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
		product := dal.GetProduct(p.Param1)
		return toJson(product)

	case "customer_price":
		// TODO extract to service
		var (
			personId  int     = p.Param1
			productId int     = p.Param2
			price     float64 = -1 // final price
		)
		if personId > 0 {
			customerPrice := dal.GetCustomerPrice(personId, productId)
			if nil != customerPrice {
				price = customerPrice.Price
			}
		}
		if price <= 0 {
			product := dal.GetProduct(p.Param2)
			if product != nil {
				price = product.Price
			}
		}
		return "json", fmt.Sprintf("{\"price\":%v}", price)
	}

	return needName()
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
