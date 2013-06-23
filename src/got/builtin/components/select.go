package components

import (
	"fmt"
	"got/core"
	// "html/template"
)

/*
   Select Component Struct

   Key is string.
   Value is string by default.

   TODO:
     support tag `param:"data"`

*/

type Select struct {
	core.Component

	// key as option value and value as label
	Data   *map[string]string // option list
	Name   string             // bind name
	Value  string             // current value/value bind
	Order  []string           // TODO: use this order(ordered map?)
	Header string             //
}

func (c *Select) Setup() {
	fmt.Println("-------------------------------select component inited.")
	// TODO auto this
	// return "template", "Select"
}

// function example
func (c *Select) IsSelected(key string) bool {
	fmt.Printf("isselected %v == %v\n", c.Value, key)
	if c.Value == key {
		return true
	}
	return false
}
