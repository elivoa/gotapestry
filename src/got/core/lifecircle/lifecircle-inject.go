package lifecircle

import (
	"fmt"
	"got/core"
	"got/debug"
	"got/utils"
	"reflect"
	"strconv"
	"strings"
)

// --------------------------------------------------------------------------------
const (
	PathInjectionTag string = "path-param"
	URLInjectionTag  string = "query" // TODO change to url
)

// injectBasic injects request&response into current Life.
func (lcc *LifeCircleControl) injectBasic() *LifeCircleControl {
	lcc.injectBasicTo(lcc.current.proton)
	return lcc
}

// InjectBasicTo will inject R & W into proton, this is not necessary, make this an option.
func (lcc *LifeCircleControl) injectBasicTo(proton core.Protoner) {
	proton.SetRequest(lcc.r)
	proton.SetResponseWriter(lcc.w)
}

func (lcc *LifeCircleControl) injectPath() *LifeCircleControl {
	lcc.injectPathTo(lcc.current.proton)
	return lcc
}

// value must be Proton struct, not ptr
// everything in lcc is belong to the root page. parameter proton is inject target.
// TODO remove reflect.
func (lcc *LifeCircleControl) injectPathTo(proton core.Protoner) {
	value := reflect.ValueOf(proton)
	t, _ := utils.RemovePointer(value.Type(), false)

	values := make(map[string][]string) // used to inject
	pathParams := extractPathParameters(lcc.r.URL.Path, lcc.pageUrl, lcc.eventName)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		var fieldKey = f.Name

		// if type is gxl.*; i.e.: a -> a.Int
		var gxlSuffix string = analysisTranslateSuffix(f.Type)
		if gxlSuffix != "" {
			fieldKey += gxlSuffix
		}

		// parse TAG: path-param | TODO Cache this.
		tagValue := f.Tag.Get(PathInjectionTag)
		if tagValue != "" {
			pathParamIndex, err := strconv.Atoi(tagValue)
			if err != nil {
				panic(fmt.Sprintf("TAG path-param must be numbers. not %v.", tagValue))
			}
			if pathParamIndex <= len(pathParams) {
				values[fieldKey] = []string{pathParams[pathParamIndex-1]}
				proton.SetInjected(f.Name, true)
			}
		}
	}
	if len(values) > 0 {
		utils.SchemaDecoder.Decode(proton, values)
	}
}

func (lcc *LifeCircleControl) injectURLParameter() *LifeCircleControl {
	lcc.injectURLParameterTo(lcc.current.proton)
	return lcc
}

func (lcc *LifeCircleControl) injectURLParameterTo(proton core.Protoner) {
	t := utils.GetRootType(proton)

	values := make(map[string][]string) // used to inject
	queries := lcc.r.URL.Query()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		var fieldKey = f.Name

		// if type is gxl.*; i.e.: a -> a.Int
		var gxlSuffix string = analysisTranslateSuffix(f.Type)
		if gxlSuffix != "" {
			fieldKey += gxlSuffix
		}

		// query param: in url query
		tagValue := f.Tag.Get(URLInjectionTag)
		if tagValue != "" {
			if tagValue == "." {
				tagValue = f.Name
			}
			if v, ok := queries[tagValue]; ok {
				values[f.Name] = v
				proton.SetInjected(f.Name, true)
				continue
			}
		}

	}
	if len(values) > 0 {
		utils.SchemaDecoder.Decode(proton, values)
	}
}
func (lcc *LifeCircleControl) injectComponentParameters(params []interface{}) *LifeCircleControl {
	lcc.injectComponentParametersTo(lcc.current.proton, params)
	return lcc
}

func (lcc *LifeCircleControl) injectComponentParametersTo(proton core.Protoner, params []interface{}) {
	// log.Printf("-621- Component [%v]'s params is: ", seg.Name)
	debugprint := false
	if debugprint {
		for i, p := range params {
			fmt.Printf("\t%3v: %v\n", i, p)
		}
		fmt.Println("\t~ END Params ~")
	}

	data := make(map[string][]string)
	var key string // key is also field name.
	for i, param := range params {
		if i%2 == 0 {
			if k, ok := param.(string); ok {
				key = fmt.Sprintf("%v%v", strings.ToUpper(k[0:1]), k[1:]) // Capitalized
				proton.SetInjected(key, true)
			} else {
				panic("component parameter must be name,value pair.")
			}
		} else {
			if key == "" || param == nil {
				// panic("value is nil")
				fmt.Println("value is nil", key, param)
				continue
			}

			// value, then set key to struct.
			switch param.(type) {
			case string:
				data[key] = []string{param.(string)}
			default: // other situation
				v := utils.GetRootValue(proton)
				injectField(v, key, param) // TO be continued....
			}
		}
	}
	if debugprint {
		debug.PrintFormMap("~ Component ~ inject component data", data)
	}

	if len(data) > 0 {
		utils.SchemaDecoder.Decode(proton, data)
	}
}

//
// inject utils
//

func injectField(target reflect.Value, fieldName string, value interface{}) {
	// target
	t := target
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// field
	field := t.FieldByName(fieldName)
	if !field.IsValid() {
		panic(fmt.Sprintf("Inject Error: Can't set '%v' to %v's%v field, type:'%v'.",
			reflect.TypeOf(value), target.Kind(), fieldName, field.Kind()))
	}
	v := reflect.ValueOf(value)

	if field.Kind() != reflect.Interface &&
		field.Kind() != reflect.Ptr &&
		v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t.FieldByName(fieldName).Set(v)
}

// *********************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************

/* ________________________________________________________________________________

   Inject Services or Parameters into proton value.
   Including:
     Request, ResponseWriter
*/
// Deprecated
// func (lcc *LifeCircleControl) InjectValue() {

// 	// 1. inject static values. (TODO: test performance)
// 	injectField(lcc.V, "W", lcc.W) // proton.W
// 	injectField(lcc.V, "R", lcc.R) // proton.R
// 	lcc.SetInjected("W", "R")      // TODO what if inject failed.

// 	// 2. inject parameter
// 	//    TODO cache tag (use map instead of loop all fields)
// 	//    How to deal with 0 and NaN, use Injected

// 	// 2.1 get value
// 	values := make(map[string][]string)
// 	t := reflect.TypeOf(lcc.Proton)
// 	if t.Kind() == reflect.Ptr {
// 		t = t.Elem()
// 	}
// 	vars := mux.Vars(lcc.R)
// 	queries := lcc.R.URL.Query()

// 	// 2.2 prepare url parameters
// 	url := lcc.R.URL.Path
// 	if !strings.HasPrefix(url, lcc.PageUrl) {
// 		panic(fmt.Sprintf("%v should has prefix %v", url, lcc.PageUrl))
// 	}

// 	// 2.3 parepare parameters
// 	paramsString := url[len(lcc.PageUrl):]
// 	if lcc.EventName != "" {
// 		index := strings.Index(paramsString, "/")
// 		if index > 0 {
// 			paramsString = paramsString[index:]
// 		}
// 	}
// 	var strParams []string
// 	if len(paramsString) > 0 {
// 		if strings.HasPrefix(paramsString, "/") {
// 			paramsString = paramsString[1:]
// 		}
// 		strParams = strings.Split(paramsString, "/")
// 	}
// 	debug.Log("-   - [injection] URL:%v, parameters:%v", url, strParams)
// 	// fmt.Printf("+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+\n")
// 	// fmt.Printf("+ url: %v\n", url)
// 	// fmt.Printf("+ lcc.PageUrl: %v\n", lcc.PageUrl)
// 	// fmt.Printf("+ paramsString: %v\n", paramsString)
// 	// fmt.Printf("+ strParams: %v\n", strParams)

// 	// ...
// 	for i := 0; i < t.NumField(); i++ {
// 		f := t.FieldByIndex([]int{i})

// 		// debug.Log("-dbg- [InjectFields] %v'th field '%v' of type %v",
// 		// 	i, f.Name, f.Type,
// 		// )

// 		// process gxl.objects
// 		var gxlSuffix string = analysisTranslateSuffix(f.Type)
// 		var fieldKey = f.Name
// 		if gxlSuffix != "" {
// 			fieldKey += gxlSuffix
// 		}

// 		var tagValue string

// 		// parse TAG: param [updated: this is not used anymore in got]
// 		tagValue = f.Tag.Get("param")
// 		if tagValue != "" {
// 			if tagValue == "." {
// 				tagValue = f.Name
// 			}
// 			v, ok := vars[tagValue]
// 			if ok {
// 				lcc.SetInjected(f.Name)
// 				values[fieldKey] = []string{v}
// 				continue
// 			}
// 		}

// 		// parse TAG: path-param
// 		tagValue = f.Tag.Get("path-param")
// 		if tagValue != "" {
// 			pathParamIndex, err := strconv.Atoi(tagValue)
// 			if err != nil {
// 				panic(fmt.Sprintf("TAG path-param must be numbers. not %v.", tagValue))
// 			}
// 			if pathParamIndex <= len(strParams) {
// 				// fmt.Printf("\t>>>>>> pathParamIndexis %v, len(strParams) = %v\n",
// 				// 	pathParamIndex, len(strParams))
// 				values[fieldKey] = []string{strParams[pathParamIndex-1]}
// 				lcc.SetInjected(f.Name)
// 			}
// 		}

// 		// query param: in url query
// 		tagValue = f.Tag.Get("query")
// 		if tagValue != "" {
// 			if tagValue == "." {
// 				tagValue = f.Name
// 			}
// 			v, ok := queries[tagValue]
// 			if ok {
// 				lcc.SetInjected(f.Name)
// 				values[f.Name] = v
// 				continue
// 			}
// 		}
// 	}
// 	if len(values) > 0 {
// 		utils.SchemaDecoder.Decode(lcc.Proton, values)
// 	}
// }

// func (lcc *LifeCircleControl) SetInjected(fields ...string) {
// 	SetInjected(lcc.v, fields...)
// 	// method := lcc.V.MethodByName("SetInjected")
// 	// if method.IsValid() {
// 	// 	for _, f := range fields {
// 	// 		method.Call([]reflect.Value{reflect.ValueOf(f), reflect.ValueOf(true)})
// 	// 	}
// 	// }
// }

// ________________________________________________________________________________
// Inject values to object, use values in lcc, but not modify any value in lcc.
// TODO: organize this, add cache of this.
//
// Deprecated
// func (lcc *LifeCircleControl) InjectValueTo(proton core.Protoner) {
// 	w, r := lcc.W, lcc.R
// 	v := reflect.ValueOf(proton)

// 	// 1. inject static values. (TODO: test performance)
// 	injectField(v, "W", w)
// 	injectField(v, "R", r)
// 	proton.SetInjected("W", true)
// 	proton.SetInjected("R", true)

// 	// 2. inject parameter
// 	// 2.1 get value
// 	values := make(map[string][]string)
// 	t, _ := utils.RemovePointer(reflect.TypeOf(proton), false)

// 	vars := mux.Vars(lcc.R)
// 	queries := r.URL.Query()

// 	// 2.2 prepare url parameters
// 	pathParams := extractPathParameters(lcc.R.URL.Path, lcc.PageUrl, lcc.EventName)

// 	// ...
// 	for i := 0; i < t.NumField(); i++ {
// 		f := t.FieldByIndex([]int{i})

// 		// debug.Log("-dbg- [InjectFields] %v'th field '%v' of type %v",
// 		// 	i, f.Name, f.Type,
// 		// )

// 		// process gxl.objects
// 		var gxlSuffix string = analysisTranslateSuffix(f.Type)
// 		var fieldKey = f.Name
// 		if gxlSuffix != "" {
// 			fieldKey += gxlSuffix
// 		}

// 		var tagValue string

// 		// parse TAG: param [updated: this is not used anymore in got]
// 		tagValue = f.Tag.Get("param")
// 		if tagValue != "" {
// 			if tagValue == "." {
// 				tagValue = f.Name
// 			}
// 			v, ok := vars[tagValue]
// 			if ok {
// 				lcc.SetInjected(f.Name)
// 				values[fieldKey] = []string{v}
// 				continue
// 			}
// 		}

// 		// parse TAG: path-param
// 		tagValue = f.Tag.Get("path-param")
// 		if tagValue != "" {
// 			pathParamIndex, err := strconv.Atoi(tagValue)
// 			if err != nil {
// 				panic(fmt.Sprintf("TAG path-param must be numbers. not %v.", tagValue))
// 			}
// 			if pathParamIndex <= len(pathParams) {
// 				values[fieldKey] = []string{pathParams[pathParamIndex-1]}
// 				lcc.SetInjected(f.Name)
// 			}
// 		}

// 		// query param: in url query
// 		tagValue = f.Tag.Get("query")
// 		if tagValue != "" {
// 			if tagValue == "." {
// 				tagValue = f.Name
// 			}
// 			v, ok := queries[tagValue]
// 			if ok {
// 				lcc.SetInjected(f.Name)
// 				values[f.Name] = v
// 				continue
// 			}
// 		}
// 	}
// 	if len(values) > 0 {
// 		utils.SchemaDecoder.Decode(proton, values)
// 	}
// }
