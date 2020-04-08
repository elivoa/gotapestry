package utils

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePinyin(t *testing.T) {
	str := ParsePinyin2("高天舒")
	strs := []string{

		"山西易露", "山西易露", "易露",
		"山西易露发", "金莲花", "Jennifer", "0", "旺华芳", "湖南赵蓓", "晒晒", "小雪", "何兴砫", "王晶妍", "小静", "DAN服饰", "卢先生", "边博华", "小薇", "章瑞庆", "吴梦洁", "水车头小兰",
	}
	for _, s := range strs {
		log.Println(">>", s, ParsePinyin2(s))
	}

	log.Println(">>", str)
}

func TestParsePinyin2(t *testing.T) {
	assert.Equal(t, ParsePinyin2("高天舒"), "gts", "计算拼音首字母错误!")
	assert.Equal(t, ParsePinyin2("李越男"), "lyn", "计算拼音首字母错误!")
	assert.Equal(t, ParsePinyin2("李玥楠"), "lyn", "计算拼音首字母错误!")
	assert.Equal(t, ParsePinyin2("ERROR"), "error", "计算拼音首字母错误!")
	// log.Println(">>", str)
}
