package main

import (
	"fmt"
	"got/utils"
	"gxl"
)

func main() {
	fmt.Println("~~~~ test utils ~~~~")
	fmt.Println(">>", utils.CurrentBasePath(), "<<")

	fmt.Println(gxl.FormatCurrency(1234567.456, 2))       //1,234,567.46
	fmt.Println(gxl.FormatCurrency(234567.456788, 3))     //234,567.457
	fmt.Println(gxl.FormatCurrency(8234567890012.456, 4)) //8,234,567,890,012.4561
	fmt.Println(gxl.FormatCurrency(0012.456, 4))          //12.4560
}
