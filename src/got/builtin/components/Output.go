/*
  Output is a Builtin Component. This component is only an component example.
  It's better to use {{.Field}} instead.
*/

package components

import (
	"fmt"
	"got/core"
)

type Output struct {
	core.Component
	Value interface{}
}

func (c *Output) Setup() (string, string) {
	fmt.Println("outout", c.Value)
	return "text", fmt.Sprint(c.Value)
}
