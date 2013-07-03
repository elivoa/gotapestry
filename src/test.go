package main

import (
	"fmt"
	"github.com/gorilla/schema"
)

type T struct {
	Ints []int
}

func main() {
	data := map[string][]string{
		"Ints": []string{"3", ""},
	}
	t := T{}
	decoder := schema.NewDecoder()
	decoder.Decode(&t, data)
	fmt.Println(t)
}
