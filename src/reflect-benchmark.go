package main

import (
	"fmt"
	"reflect"
	"time"
)

func main() {

	testReflect(1000000)
}

// TOBE continued.
func testReflect(n int64) {
	var count int64 = 1 // targetp

	// --------------------------------------------------------------------------------
	s := time.Now().UnixNano()
	var i int64
	for i = 0; i < n; i++ {
		count = i * i
	}
	base := time.Now().UnixNano() - s
	fmt.Printf("[%10v] %v\n", base, "// base directly assign")

	s2 := time.Now().UnixNano()
	var o int64
	for o = 0; o < n; o++ {
		v := reflect.ValueOf(&count)
		v = v.Elem()
		v.SetInt(o)
	}
	time2 := time.Now().UnixNano() - s2
	fmt.Printf("[%10v] (%v) %v\n", time2, time2/base, "// reflect.value and assing.")

	s3 := time.Now().UnixNano()
	v := reflect.ValueOf(&count)
	v = v.Elem()
	for o = 0; o < n; o++ {
		v.SetInt(o)
	}
	time3 := time.Now().UnixNano() - s3
	fmt.Printf("[%10v] (%v) %v\n", time3, time3/base, "// set use reflect")
}
