package components

import (
	"fmt"
	"got/core"
)

/* _____________________________________________________
   Component Object
*/
type ProvinceSelect struct {
	core.Component

	Data   *map[string]string // option list
	Name   string
	Value  string
	Header string //
}

func (c *ProvinceSelect) Setup() {
	fmt.Println("------------------------------- Province Select component inited.")
	c.Data = &provinceData
	if !c.Injected("Header") {
		c.Header = "省"
	}
}

// function example
func (c *ProvinceSelect) IsSelected(key string) bool {
	fmt.Printf("isselected %v == %v\n", c.Value, key)
	if c.Value == key {
		return true
	}
	return false
}

var (
	provinceData = map[string]string{
		//	"--- 直辖市 ---": "--- 直辖市 ---",
		"北京市": "北京市",
		"天津市": "天津市",
		"上海市": "上海市",
		"重慶市": "重慶市",

		//	"--- 省 ---": "--- 省 ---",
		"安徽省": "安徽省",
		"福建省": "福建省",
		"甘肃省": "甘肃省",
		"广东省": "广东省",
		"贵州省": "贵州省",
		"河北省": "河北省",
		"龙江省": "龙江省",
		"河南省": "河南省",
		"湖北省": "湖北省",
		"湖南省": "湖南省",
		"吉林省": "吉林省",
		"江西省": "江西省",
		"江苏省": "江苏省",
		"辽宁省": "辽宁省",
		"山东省": "山东省",
		"陕西省": "陕西省",
		"山西省": "山西省",
		"四川省": "四川省",
		"云南省": "云南省",
		"浙江省": "浙江省",
		"青海省": "青海省",
		"海南省": "海南省",
		"台湾":  "台湾",

		//	"--- 自治区 ---": "--- 自治区 ---",
		"广西壮族自治区":  "广西壮族自治区",
		"内蒙古自治区":   "内蒙古自治区",
		"宁夏回族自治区":  "宁夏回族自治区",
		"西藏自治区":    "西藏自治区",
		"新疆维吾尔自治区": "新疆维吾尔自治区",
		"香港特别行政区":  "香港特别行政区",
		"澳門特别行政区":  "澳門特别行政区",
	}
)
