package main

import (
	"fmt"
	"time"
)

func main() {
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
