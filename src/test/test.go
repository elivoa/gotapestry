package main

import (
	"fmt"
	"syd/dal/statdao"
	"syd/service/orderservice"
	// "syd/service/personservice"
	"strconv"
	"strings"
)

var va = 1 << 0
var vb = 1 << 1
var vc = 1 << 2
var vd = 1 << 3

func main() {

	a := "http://www.chanel.com/dam/fashion/catalog/collections/15K/RTW/looks/15K####.jpg.fashionImg.veryhi.jpg"
	for i := 50; i < 100; i++ {
		// b := strconv.Itoa( + i)[1:5]
		// fmt.Println(b)
		fmt.Println(strings.Replace(a, "####", strconv.Itoa(i), -1))
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
