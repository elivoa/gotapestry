package gxl

import (
	"bytes"
	"fmt"
	"strconv"
)

func FormatCurrency(c float64, digit int) string {
	str := fmt.Sprintf("%."+strconv.Itoa(digit)+"f", c)
	leading := len(str) - digit - 1
	n := leading / 3
	r := leading % 3
	var result bytes.Buffer
	if c < 0 { // sign
		result.WriteString("-")
	}
	if r > 0 {
		result.WriteString(str[0:r])
	}
	for i := 0; i < n; i++ {
		if (i == 0 && r > 0) || i > 0 {
			result.WriteString(",")
		}
		result.WriteString(str[r+i*3 : r+(i+1)*3])
	}
	result.WriteString(str[leading:])
	return result.String()
}
