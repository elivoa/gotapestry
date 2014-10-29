package main

import (
	"fmt"
	"github.com/elivoa/got/utils"
	"syd/dal/statdao"
	"syd/service/orderservice"
	// "syd/service/personservice"
	"syd/dal/persondao"
	"syd/service/statservice"
	"time"
)

var va = 1 << 0
var vb = 1 << 1
var vc = 1 << 2
var vd = 1 << 3

func main() {
	// now := time.Now()
	// end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0)
	// start := end.AddDate(0, 0, 1)
	a := va | vc
	fmt.Println("va ", a&va > 0)
	fmt.Println("vb ", a&vb > 0)
	fmt.Println("vc ", a&vc > 0)
	fmt.Println("vd ", a&vd > 0)
	// userdao.ListUserByIdSet(1, 3, 4)

	fmt.Println("----------------------------------------------------------------------------------------------------")
	fmt.Println(time.Now())
	var t time.Time
	fmt.Println(t)
	fmt.Println(utils.ValidTime(time.Now()))
	// for i := 0; i < 20; i++ {
	// 	testLoad3()
	// 	time.Sleep(time.Microsecond * 500)
	// }

	// fmt.Println("first call done!")

	// fmt.Println("Waiting 2 seconds...")
	// time.Sleep(time.Second * 8)
}

func testLoad3() {
	// fmt.Println("---- start loading...")
	// customer := personservice.GetPerson(25)
	customer, err := persondao.Get(35)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("-----------------", customer)
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
