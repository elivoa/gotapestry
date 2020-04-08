package main

import (
	"fmt"
	"syd"

	// _ "syd/generated"
	_ "syd"
	_ "syd/dal/userdao"
	_ "syd/model"
	_ "syd/service"

	_ "github.com/elivoa/got"
	_ "github.com/elivoa/gxl"
)

// TESTing....
func main() {
	fmt.Println("TODO: Make this project a github.com project.")
	syd.Start()
}
