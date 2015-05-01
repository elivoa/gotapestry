package main

import (
	"fmt"
	"github.com/elivoa/got/builtin/services"
	"time"
)

func main() {
	fmt.Println("------------------------------------------------------------------------------------------")
	var parameters = map[string]interface{}{
		"provider": 33,
		"from":     time.Now(),
	}

	fmt.Println(services.Link.GeneratePageUrlWithContextAndQueryParameters("inventory", parameters))
	fmt.Println("------------------------------------------------------------------------------------------")

}
