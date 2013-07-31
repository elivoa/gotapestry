package main

import (
	"fmt"
	"syd/service/suggest"
)

func main() {
	suggest.EnsureLoaded()
	items, err := suggest.Lookup("w", "customer")
	if err != nil {
		panic(err.Error())
	}

	for i, v := range items {
		fmt.Printf("%v: %v\n", i, v)
	}

}
