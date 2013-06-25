package templates

import (
	// "encoding/json"
	"fmt"
	"html/template"
	"log"
	"strings"
)

/*_______________________________________________________________________________
  Register components
  TODO:
    Register outside, make this a framework.
    . generate select component
*/
func registerComponentFuncs() {
	// init functions
	Templates.Funcs(template.FuncMap{})
}

func RegisterComponent(name string, f interface{}) {
	funcName := fmt.Sprintf("t_%v", strings.Replace(name, "/", "_", -1))
	debuglog("-108- [RegisterComponent] ", funcName)
	Templates.Funcs(template.FuncMap{funcName: f})
}

// --------------------------------------------------------------------------------
// log
//
var debugLog = true

func debuglog(format string, params ...interface{}) {
	if debugLog {
		log.Printf(format, params...)
	}
}

/*_______________________________________________________________________________
  Test
*/

// type SelectModel struct {
// 	Data  string
// 	Value string
// 	X     int
// 	B     bool
// }

// func SelectComponentTest(params ...interface{}) string {
// 	// 1. parse parameters.
// 	// .. support one string or many strings. splited by ,
// 	var model = SelectModel{}
// 	Unmarshal(&model, params...)
// 	fmt.Printf(">>>>> %v <<<<<\n", model)
// 	// utils.SchemaDecoder.Decode(&model, *data)
// 	return fmt.Sprintf("[> %v <]", model)
// }

/*
   Unmarshal parameters into map

   Support style in template file:
     {{t_select "data:ii" "value:9" "x:3888"}}

   TODO:
     support [space , comma, ...] in string
     support multi param in one parameter.
*/

/* ------------------------------------------------------------ */

// // this is test version, modify this into got style.
// func SelectComponent(parameters ...string) string {
// 	// 1. parse parameters.
// 	// .. support one string or many strings. splited by ,

// 	var model = SelectModel{}
// 	data := Unmarshal(parameters...)
// 	utils.SchemaDecoder.Decode(&model, *data)
// 	return fmt.Sprintf("[> %v <]", model)
// }

// /*
//    Unmarshal parameters into map, to use by gorilla/schema

//    Support style in template file:
//      {{t_select "data:ii" "value:9" "x:3888"}}

//    TODO:
//      support [space , comma, ...] in string
//      support multi param in one parameter.
// */
// func Unmarshal(parameters ...string) *map[string][]string {
// 	params := make(map[string][]string)
// 	for _, piece := range parameters {
// 		sps := strings.Split(piece, ":")
// 		var key string // capitalized key, because unexported field has no meaning.
// 		if len(sps[0]) >= 1 {
// 			key = fmt.Sprintf("%v%v", strings.ToUpper(sps[0][0:1]), sps[0][1:])
// 			fmt.Printf("key is  %v\n", key)
// 		}
// 		if len(sps) == 1 {
// 			params[key] = []string{"true"}
// 		} else if len(sps) == 2 {
// 			params[key] = []string{sps[1]}
// 		}
// 	}
// 	return &params
// }

// func Test(str string) string {
// 	//parse.
// 	//Templates.AddParseTree
// 	fmt.Println()
// 	return ""
// }
