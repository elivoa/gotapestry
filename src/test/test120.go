package main

import (
	"fmt"
	"github.com/elivoa/gxl"
	"time"
)

func main() {
	now := time.Now().UTC()

	fmt.Println("now is : ", now)
	s, e := gxl.NatureTimeRangeUTC(0, 0, 0)
	fmt.Println("debug:::: ", s, " ////  ", e)
	fmt.Println("----------------------------")
	t, err := time.ParseInLocation("2006-01-02 15:04:05", "2015-01-01 01:33:22", time.Local)
	if err != nil {
		panic(err)
	}
	fmt.Println("+8now is:", t)

	start := now.Truncate(time.Hour * 24)
	end := now.AddDate(0, 0, 1).Truncate(time.Hour * 24)
	fmt.Println(start)
	fmt.Println(end)

	fmt.Println(time.Now().Format("2006-01-02"))
	fmt.Println(time.Now().Truncate(time.Hour * 24))

	fmt.Println("--", time.Now().AddDate(0, 0, 1).Format("2006-01-02"))
	ss := datekeys(7)
	for _, k := range ss {
		fmt.Println("LL ", k)
	}
}

func datekeys(lastNDays int) []string {
	t := time.Now().AddDate(0, 0, -lastNDays+1)
	result := []string{}
	for i := 0; i < lastNDays; i++ {
		result = append(result, t.AddDate(0, 0, i).Format("2006-01-02"))
	}

	return result
}
