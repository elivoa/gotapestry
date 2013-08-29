package main

import (
	"fmt"
	"got/utils"
	"syd/service/statservice"
	"time"
)

func main() {
	now := time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0)
	start := end.AddDate(0, 0, 1)
}

func main2() {
	t := utils.NewTimer()
	// end := time.Now()
	// start := end.AddDate(0, 0, -7)
	dd := statservice.CalcHotSaleProducts(0, 0, -7)
	fmt.Println(dd)
	// list, err := statservice.HotSaleProducts(0, 0, 7)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// fmt.Println("total count", len(list))
	// fmt.Println(list)

	fmt.Println("execution time is: ", t.NowSecond())
}
