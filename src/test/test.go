package main

import (
	"fmt"
	"got/utils"
	"gxl"
	"regexp"
)

func main() {
	t := utils.NewTimer()

	fmt.Println(gxl.FormatCurrency(-12384.333, 2))

	fmt.Println("~~~~ test utils ~~~~")
	fmt.Println(">>", utils.CurrentBasePath(), "<<")

	fmt.Println(gxl.FormatCurrency(1234567, 0))       //1,234,567.46
	fmt.Println(gxl.FormatCurrency(234567.456788, 3))     //234,567.457
	fmt.Println(gxl.FormatCurrency(8234567890012.456, 4)) //8,234,567,890,012.4561
	fmt.Println(gxl.FormatCurrency(0012.456, 4))          //12.4560
	fmt.Println("********************************************************************************")

	var rePrintValue, _ = regexp.Compile("^(.*){{(.*)}}$")
	result := rePrintValue.FindStringSubmatch("/order/list/{{.IDIDID}}")
	for _, r := range result {
		fmt.Println(r)
	}

	fmt.Println("execution time is: ", t.NowSecond())
}
