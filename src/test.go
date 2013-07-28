package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"got/templates/transform"
	"io"
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

	// trans := transform.NewTransformer()
	// trans.Parse(r)
	// str := trans.RenderToString()
	// // str, _ := transform.TransformTemplate(r)
	// fmt.Println(str)
	parseTemplate(r)
}

func parseTemplate(reader io.Reader) {
	decoder := xml.NewDecoder(reader)
	decoder.DecodeElement(nil, nil)
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
