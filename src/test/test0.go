package main

import (
	"fmt"
	// "syd/service/personservice"
	"bytes"
	"strings"
)

func main() {
	var str = "   test this string     \n    sdfj \n"
	var bytes = []byte(str)
	fmt.Println("FORMAT: ", str)
	fmt.Println("AGTER : ", string(TrimTextNode(bytes)))
	fmt.Println(strings.Replace("aslkdjfaldsjf \n lskdjflsdj", "\n", "\\n", -1))
}

// TODO trim node function not finished.
func TrimTextNode(text []byte) []byte {
	fmt.Printf(">>[%s]\n", string(text))

	var (
		firstValidCharacter int  = 0
		hasUsefulCharacters bool = false
	)

	for _, b := range text {
		switch b {
		case ' ':
			// pass
		case '\r', 'n':

		default: // has other characters
			hasUsefulCharacters = true
		}
	}
	fmt.Println(firstValidCharacter, hasUsefulCharacters)
	// var (
	// 	addSpaceLeft    bool = false
	// 	addSpaceRight   bool = false
	// 	addNewLineLeft  bool = false
	// 	addNewLineRight bool = false
	// )
	// for _, b := range bytes {
	// 	// if b
	// }
	return bytes.Trim(text, " \r\n")
}
