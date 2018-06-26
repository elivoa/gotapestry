package utils

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePinyin(t *testing.T) {
	str := ParsePinyin("高天舒")

	log.Println(">>", str)
}

func TestParsePinyin2(t *testing.T) {
	assert.Equal(t, ParsePinyin2("高天舒"), "gts", "计算拼音首字母错误!")
	assert.Equal(t, ParsePinyin2("李越男"), "lyn", "计算拼音首字母错误!")
	assert.Equal(t, ParsePinyin2("李玥楠"), "lyn", "计算拼音首字母错误!")
	assert.Equal(t, ParsePinyin2("ERROR"), "error", "计算拼音首字母错误!")
	// log.Println(">>", str)
}
