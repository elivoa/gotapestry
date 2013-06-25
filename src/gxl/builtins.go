/*
  Go Extend Library
*/
package gxl

import ()

/* ________________________________________________________________________________
   Int values
*/

func NewInt(i int) *Int {
	return &Int{i}
}

type Int struct {
	Int int
}

// func (i *Int) Int() int {
// 	return Int
// }

func (i *Int) String() string {
	return string(i.Int)
}

func (i *Int) Set(value int) *Int {
	i.Int = value
	// i.initialized = true
	return i
}

// func (i *Int) IsSet() bool {
// 	return i.initialized
// }

// func (i *Int) String() string {
// 	if i.initialized {
// 		return "int{nil}"
// 	}
// 	return string(i.int)
// }
