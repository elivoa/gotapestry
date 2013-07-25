package main

import (
	"fmt"
)

type OrderType uint

func (t OrderType) Coercion(from string) OrderType {
	return t
}

func main() {
	var a OrderType
	fmt.Println(a)
	var b uint = 3
	a = uint(b)
	fmt.Println(a)
	// fmt.Println(Wholesale)
	// fmt.Println(Wholesale.Coercion("44"))
	// fmt.Println(Wholesale.Coercion("44"))
}
