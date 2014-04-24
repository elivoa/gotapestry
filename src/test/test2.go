package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/elivoa/got/templates/transform"
	"reflect"
	"strings"
	"syd/exceptions"
)

var html = `
<html>
  <head>
    <meta a="b"></meta>
    <link href="/static/bocss" rel="stylesheet" type="text/css">====
    <title>haha</title>
    <meta>+++++++++++++
  </head>
  <body>
    <t:block id="XXX">
      <h1>CheDan
    </t:block>
  </body>
</html>
`

func main() {

	fmt.Printf("reflect of error %s\n", reflect.TypeOf(errors.New("DDD")))

	fmt.Println("--------------------------------------------------------------------------------")

	t := transform.NewTransformer()
	t.Parse(strings.NewReader(html))
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println(t.RenderToString())

	abc := sha256.Sum256([]byte("elivoa"))
	ddd := sha1.Sum([]byte("elivoa"))
	fmt.Println("\n", string(abc[:]))
	fmt.Println("\n", string(ddd[:]))

	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case error:
				fmt.Println("type is error")
			case string:
				fmt.Println("type is string")
			case exceptions.LoginError:
				fmt.Println("type is LoginError")
			}
		}
	}()
	err := TestSome()
	print(err)

}
func TestSome() error {
	panic(&exceptions.LoginError{Message: "slkdfjalsdjflkj"})
}
