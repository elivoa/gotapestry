package main

import (
	"fmt"
	"github.com/elivoa/got/utils"
	"syd/dal/statdao"
	"syd/service/orderservice"
	// "syd/service/personservice"
	"strconv"
	"strings"
	"syd/service/statservice"
)

var va = 1 << 0
var vb = 1 << 1
var vc = 1 << 2
var vd = 1 << 3

func main() {

	a := "http://tittyandco.net/fw/shop/img/product/TIT14F####/TIT14F####_pz_a001.jpg"
	for i := 301; i < 350; i++ {
		b := strconv.Itoa(10000 + i)[1:5]
		// fmt.Println(b)
		fmt.Println(strings.Replace(a, "####", b, -1))
	}

}

func testLoad() {
	fmt.Println("---- start loading...")
	stats, err := statdao.TodayStat(20)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(stats)
	if len(stats) > 0 {
		fmt.Println("-----------------")
		for i, s := range stats {
			fmt.Println(i, s.Id, s.NOrder, s.TotalPrice)
		}
	}
}

func testLoad2() {
	fmt.Println("---- start loading...")
	orders, err := orderservice.ListOrder("all")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(len(orders))
	fmt.Println("-----------------", len(orders))
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
