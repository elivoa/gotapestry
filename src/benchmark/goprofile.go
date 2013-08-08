package main

import (
	"flag"
	"log"
	"os"
	"reflect"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

type T struct {
	A []string
	B []int
	C float64
}

func main() {
	test()
	log.Println("all test done!")
}
func test() {

	f, _ := os.Create("ls.prof")

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer func() {
			pprof.StopCPUProfile()
		}()
	}

	// ...
	t := &T{}
	for i := 0; i < 1000000; i++ {
		v := reflect.ValueOf(t)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		t.B = append(t.B, i)
		t.A = append(t.A, sumtest())
	}

	// var buffer bytes.Buffer
	pprof.WriteHeapProfile(f)
	// fmt.Println(buffer.String())

}
func sumtest() string {
	return "abc"
}
