package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("------------------------------------------------------------------------------------------")
	fmt.Println(time.Now())
	fmt.Println(time.Now().Truncate(time.Hour * 24))
}
