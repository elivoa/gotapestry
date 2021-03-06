package components

import (
	"github.com/elivoa/got/core"
	"syd/service"
	"syd/service/personservice"
)

/*
   This component generate 2 html elements.
   1. input:hidden, id=ClientId_id
      save id of the element.
      will submit id to form.
   2. input:text, id=ClientId_query
      accept query input.
      display final text.
*/
type SuggestControl struct {
	core.Component

	ClientIdd    string // client id, in html.
	Name         string // form submit key
	Value        int    // value. should be id in this case
	DisplayValue string // value to show
	Category     string // category of select. TODO support ALL
	Callback     string // javascript callback method ?no-use
	MultiMode    bool   // if multiline mode, it's more completed.
}

func (c *SuggestControl) Setup() { // (string, string) {
	if !c.Injected("ClientIdd") {
		c.ClientIdd = "factory_suggest"
	}
	c.initSuggest()
	// return "template", "SuggestControl"
}

func (c *SuggestControl) initSuggest() {
	// id, err := strconv.Atoi(c.Value)
	id := c.Value
	switch c.Category {
	case "factory":
		person, err := service.Person.GetPersonById(id)
		if err != nil {
			panic(err)
		}
		if person != nil {
			c.DisplayValue = person.Name
		}
	case "customer":
		person := personservice.GetProducer(id)
		if person != nil {
			c.DisplayValue = person.Name
		}
	case "product":
		product, err := service.Product.GetFullProduct(id)
		if err == nil && product != nil {
			c.DisplayValue = product.Name
		}
	default:
	}
}
