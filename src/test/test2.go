package main

import (
	"fmt"
	"got/templates/transform"
	"strings"
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
	t := transform.NewTransformer()
	t.Parse(strings.NewReader(html))
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println(t.RenderToString())
}
