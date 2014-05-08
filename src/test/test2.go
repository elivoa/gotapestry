package main

import (
	"bytes"
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

func test(out *bytes.Buffer, i int) {
	out.WriteString("(-)")
	if i > 0 {
		i = i - 1
		test(out, i)
	}
}

// fu

func main() {
	treestr := "(  d  )"
	treestr = strings.Replace(treestr, " ", "&nbsp;", -1)
	print(treestr)

}
