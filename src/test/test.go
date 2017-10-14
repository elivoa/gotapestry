package main

import (
	"fmt"
	"syd/dal/statdao"
	"syd/service/orderservice"
	// "syd/service/personservice"
	"bytes"
	"strconv"
	"strings"
	"time"
)

var va = 1 << 0
var vb = 1 << 1
var vc = 1 << 2
var vd = 1 << 3

// 例如 0 0 -1, 返回最近两天的时间点。
func NatureTimeRangeUTC(years, months, days int) (start, end time.Time) {
	Timezone := 0
	natureEnd := time.Now().AddDate(0, 0, 1).UTC().Truncate(time.Hour * 24).
		Add(time.Hour * time.Duration(-Timezone))
	natureStart := natureEnd.AddDate(years, months, days-1)
	return natureStart, natureEnd
}

func EndOfTodayUTC() (t time.Time) {
	year, month, day := time.Now().Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, time.UTC)
}
func StartOfTomorrowUTC() (t time.Time) {
	year, month, day := time.Now().AddDate(0, 0, 1).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func UntilEndOfTodayRangeUTC(days int) (start, end time.Time) {
	year, month, day := time.Now().AddDate(0, 0, -(days - 1)).Date()
	startTime := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return startTime, EndOfTodayUTC()
}

func UntilStartOfTomorrowRangeUTC(days int) (start, end time.Time) {
	year, month, day := time.Now().AddDate(0, 0, -(days - 1)).Date()
	startTime := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return startTime, StartOfTomorrowUTC()
}

func main() {
	Timezone := 0

	// start, end := NatureTimeRangeUTC(0, 0, 0)
	// fmt.Println(start, end)
	fmt.Println("Now is : ", time.Now())
	fmt.Println(UntilStartOfTomorrowRangeUTC(1))
	fmt.Println(UntilStartOfTomorrowRangeUTC(2))
	fmt.Println(UntilStartOfTomorrowRangeUTC(3))
	fmt.Println(UntilStartOfTomorrowRangeUTC(4))
	fmt.Println(UntilStartOfTomorrowRangeUTC(5))
	fmt.Println(UntilStartOfTomorrowRangeUTC(6))
	fmt.Println(time.Now().AddDate(0, 0, 1))
	fmt.Println(time.Now().AddDate(0, 0, 1).UTC())
	fmt.Println(time.Now().AddDate(0, 0, 1).UTC().Truncate(time.Hour * 24))
	fmt.Println(time.Now().AddDate(0, 0, 1).UTC().Truncate(time.Hour * 24).Add(time.Hour * time.Duration(-Timezone)))
	fmt.Println("========================================")

	fmt.Println(00011 &^ 0110)
	var buf bytes.Buffer
	var PLACEHOLDER string = "(____PageHeadBootstrap_replace_to_html____)"
	buf.WriteString("something  b((((((((()efore ((____PageHeadBootstrap_replace_to_html____) end")
	fmt.Println("before: ", buf.String())
	newBuf := replaceInBuffer(&buf, PLACEHOLDER, "-replace to this-")
	fmt.Println("after : ", newBuf.String())
}

func replaceInBuffer(buf *bytes.Buffer, s, to string) *bytes.Buffer {
	var newbuf bytes.Buffer
	var err error
	var c byte
	var position = 0
	var cursor = 0
	var found_first int
	var matchedbytes = make([]byte, len(s))
	for err == nil {
		if c, err = buf.ReadByte(); err != nil {
			break // eof
		}
		// fmt.Printf("outer: char:%s position:%d\n", string(c), position)

		// fmt.Println("c == s[cursor]:  ", string(c), "==", string(s[cursor]), " cursor: ", cursor)
		if c == s[cursor] || c == s[0] {
			// fmt.Printf("inner: char:%s position:%d cursor:%d\n", string(c), position, cursor)
			if cursor == 0 {
				found_first = position
			} else if c == s[0] {
				// special   中途直接重启match e.g.: ((
				// fmt.Println(">> c == s[0]:  ", string(c), "==", string(s[cursor]), " cursor: ", cursor)
				newbuf.Write(matchedbytes[:cursor])
				cursor = 0
				found_first = position
			}
			// fmt.Println("cursor == len(s)-1:  ", cursor, "=?", len(s))
			if cursor == len(s)-1 { // found match
				// fmt.Println("Found cursor is: ", cursor, " char is ", string(s[cursor]))
				newbuf.WriteString(to)
				break // found
			}
			matchedbytes[cursor] = c
			// fmt.Println("matchedbytes:", string(matchedbytes))
			cursor += 1
		} else if found_first > 0 {
			// find part, start failed. reset.
			newbuf.Write(matchedbytes[:cursor])
			found_first = 0
			cursor = 0
		}

		if found_first == 0 { // find part, reset
			newbuf.WriteByte(c)
		}
		position += 1
	}
	newbuf.Write(buf.Bytes())
	return &newbuf
}

func main2() {

	a := "http://www.chanel.com/dam/fashion/catalog/collections/15K/RTW/looks/15K####.jpg.fashionImg.veryhi.jpg"
	for i := 50; i < 100; i++ {
		// b := strconv.Itoa( + i)[1:5]
		// fmt.Println(b)
		fmt.Println(strings.Replace(a, "####", strconv.Itoa(i), -1))
	}

}

func testLoad() {
	fmt.Println("---- start loading...")
	stats, err := statdao.TodayStat(time.Now(), 20)
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
