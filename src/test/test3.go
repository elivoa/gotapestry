package main

import (
	"fmt"
	"github.com/elivoa/got/builtin/services"
	"log"
	"regexp"
	"strings"
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
	t := time.Now()
	t2 := time.Now()
	fmt.Println(t.UnixNano() < t2.UnixNano())
	fmt.Println(t.UnixNano() - t2.UnixNano())
	fmt.Printf("%s", t.Format("2016-01-02 06:06:07"))

	s := "season=2015-08-09"

	reg, err := regexp.Compile("[^A-Za-z0-9-]+")
	if err != nil {
		log.Fatal(err)
	}

	safe := reg.ReplaceAllString("a*-+fe5v9034,j*.AE6", "-")
	safe = strings.ToLower(strings.Trim(safe, "-"))
	fmt.Println(safe) // Output: a*-+fe5v9034,j*.ae6
	fmt.Println(reg.ReplaceAllString(s, ""))

}
