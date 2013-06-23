package main

import (
	// "encoding/json"
	"fmt"
	// "reflect"
	"math/rand"
	// "syd/model"
	"syd/service/suggest"
	"time"
)

type Protoner interface {
	SXXX()
}

type Pager interface {
	Protoner
}

type Proton struct {
	a string
}

func (p *Proton) SXXX() {
	fmt.Println("aaa")
}

type Page struct {
	Proton
}

/* ********************************** */
func test(p interface{}) Protoner {
	return p.(Protoner)
}

func main2() {
	p := Page{}
	fmt.Printf(">>>>>>>>>>>>>> %v\n", p)
	newp := test(&p)
	fmt.Printf(">>>>>>>>>>>>>> %v\n", newp)
	//	fmt.Println(test(p))
	// fmt.Println(test(p).Request())

	fmt.Println(p)

	fmt.Println("---------------------------------------")
	suggest.Load()
	suggest.PrintAll()
	items, err := suggest.Lookup("s", "customer")
	fmt.Println("------------------------------------------------------")
	if err == nil {
		for i, item := range *items {
			fmt.Printf("\t%v: %v | %v\n", i, item.Text, item.QuickString)
		}
	} else {
		panic(err)
	}

}

func main3() {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println(88130430888)
	fmt.Println(13052988992)

	a := fmt.Sprintf("%v%v", time.Now().Format("060102030405"), rand.Intn(999))
	fmt.Println(a)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

type Integer struct {
	int
	initialized bool
}

func Int(i int) Integer {
	return Integer{i, true}
}

func main() {
	a := Int(8)
	fmt.Println(a)
	fmt.Println(a.int + 8)
}
