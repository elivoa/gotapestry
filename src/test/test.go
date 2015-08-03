package main

import (
	"fmt"
	"syd/dal/statdao"
	"syd/service/orderservice"
	// "syd/service/personservice"
	"bytes"
	"strconv"
	"strings"
)

var va = 1 << 0
var vb = 1 << 1
var vc = 1 << 2
var vd = 1 << 3

func main() {
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
