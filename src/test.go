package main

import (
	"fmt"
)

func main() {
	a := []int{1, 2, 3, 4}
	fmt.Println(a[0 : len(a)-2])
	fmt.Println(a[len(a)-1])
	
}
