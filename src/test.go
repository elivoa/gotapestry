package main

import (
	"bufio"
	"code.google.com/p/go.net/html"
	"fmt"
	"got/templates/transform"
	"os"
)

func main() {
	// open input file
	fi, err := os.Open("test.html")
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	fmt.Println("********************************************************************************")
	// make a read buffer
	r := bufio.NewReader(fi)
	// read source from file.

	// b, err := ioutil.ReadFile("test.html")
	// if err != nil {
	// 	panic("error")
	// }

	// str := string(b)
	// fmt.Println(html)
	// trans := transform.NewTransformer()
	// trans.Parse(r)
	// str := trans.RenderToString()
	// // str, _ := transform.TransformTemplate(r)
	// fmt.Println(str)
	// ParseTemplate(html)
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if z.Err().Error() == "EOF" {
				return
			} else {
				panic(z.Err().Error())
			}
			// ...
			fmt.Println("return")
			return
		}
		// Process the current token.
		fmt.Println(">>       ", tt)
		fmt.Println("Text:    ", string(z.Text()))
		// fmt.Println("TagName: ", string(z.TagName()))
		// fmt.Println("TagAttr: ", string(z.TagAttr()))
	}

	fmt.Println(z)

}

// --------------------------------------------------------------------------------
func main222() {
	// open input file
	fi, err := os.Open("test.html")
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	fmt.Println("********************************************************************************")
	// make a read buffer
	r := bufio.NewReader(fi)

	trans := transform.NewTransformer()
	trans.Parse(r)
	str := trans.RenderToString()
	// str, _ := transform.TransformTemplate(r)
	fmt.Println(str)
}
