package lifecircle

import (
	"fmt"
	"github.com/gorilla/mux"
	"got/debug"
	"got/utils"
	"reflect"
	"strconv"
	"strings"
)

/* ________________________________________________________________________________
   Inject Services or Parameters into proton value.
   Including:
     Request, ResponseWriter
*/
func (lcc *LifeCircleControl) InjectValue() {

	// fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	// fmt.Println(reflect.TypeOf(lcc.Proton))
	// fmt.Printf(">>>> %v\n", reflect.TypeOf(lcc.W))
	// fmt.Printf(">>>> %v\n", reflect.TypeOf(lcc.R))
	// fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

	// 1. inject static values. (TODO: test performance)
	injectField(lcc.V, "W", lcc.W) // proton.W
	injectField(lcc.V, "R", lcc.R) // proton.R
	lcc.SetInjected("W", "R")      // TODO what if inject failed.

	// 2. inject parameter
	//    TODO cache tag (use map instead of loop all fields)
	//    How to deal with 0 and NaN, use Injected

	// 2.1 get value
	values := make(map[string][]string)
	t := reflect.TypeOf(lcc.Proton)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	vars := mux.Vars(lcc.R)
	queries := lcc.R.URL.Query()

	// 2.2 prepare url parameters
	url := lcc.R.URL.Path
	if !strings.HasPrefix(url, lcc.PageUrl) {
		panic(fmt.Sprintf("%v should has prefix %v", url, lcc.PageUrl))
	}

	// 2.3 parepare parameters
	paramsString := url[len(lcc.PageUrl):]
	if lcc.EventName != "" {
		index := strings.Index(paramsString, "/")
		if index > 0 {
			paramsString = paramsString[index:]
		}
	}
	var strParams []string
	if len(paramsString) > 0 {
		if strings.HasPrefix(paramsString, "/") {
			paramsString = paramsString[1:]
		}
		strParams = strings.Split(paramsString, "/")
	}
	debug.Log("-   - [injection] URL:%v, parameters:%v", url, strParams)
	// fmt.Printf("+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+\n")
	// fmt.Printf("+ url: %v\n", url)
	// fmt.Printf("+ lcc.PageUrl: %v\n", lcc.PageUrl)
	// fmt.Printf("+ paramsString: %v\n", paramsString)
	// fmt.Printf("+ strParams: %v\n", strParams)

	// ...
	for i := 0; i < t.NumField(); i++ {
		f := t.FieldByIndex([]int{i})

		// debug.Log("-dbg- [InjectFields] %v'th field '%v' of type %v",
		// 	i, f.Name, f.Type,
		// )

		// process gxl.objects
		var gxlSuffix string = analysisTranslateSuffix(f.Type)
		var fieldKey = f.Name
		if gxlSuffix != "" {
			fieldKey += gxlSuffix
		}

		var tagValue string

		// parse TAG: param [updated: this is not used anymore in got]
		tagValue = f.Tag.Get("param")
		if tagValue != "" {
			if tagValue == "." {
				tagValue = f.Name
			}
			v, ok := vars[tagValue]
			if ok {
				lcc.SetInjected(f.Name)
				values[fieldKey] = []string{v}
				continue
			}
		}

		// parse TAG: path-param
		tagValue = f.Tag.Get("path-param")
		if tagValue != "" {
			pathParamIndex, err := strconv.Atoi(tagValue)
			if err != nil {
				panic(fmt.Sprintf("TAG path-param must be numbers. not %v.", tagValue))
			}
			if pathParamIndex <= len(strParams) {
				// fmt.Printf("\t>>>>>> pathParamIndexis %v, len(strParams) = %v\n",
				// 	pathParamIndex, len(strParams))
				values[fieldKey] = []string{strParams[pathParamIndex-1]}
				lcc.SetInjected(f.Name)
			}
		}

		// query param: in url query
		tagValue = f.Tag.Get("query")
		if tagValue != "" {
			if tagValue == "." {
				tagValue = f.Name
			}
			v, ok := queries[tagValue]
			if ok {
				lcc.SetInjected(f.Name)
				values[f.Name] = v
				continue
			}
		}
	}
	if len(values) > 0 {
		utils.SchemaDecoder.Decode(lcc.Proton, values)
	}
}

func analysisTranslateSuffix(t reflect.Type) string {
	switch t.String() {
	case "*gxl.Int":
		return ".Int"
	}
	return ""
}

func (lcc *LifeCircleControl) InjectComponentParameters(params []interface{}) *LifeCircleControl {
	// debug log .....
	// log.Printf("-621- Component [%v]'s params is: ", seg.Name)
	debug := false
	if debug {
		for i, p := range params {
			fmt.Printf("\t%3v: %v\n", i, p)
		}
		fmt.Println("\t~ END Params ~")
	}

	data := make(map[string][]string)
	var key string
	for i, param := range params {
		if i%2 == 0 {
			var ok bool
			key, ok = param.(string)
			if !ok {
				panic("component parameter must be name,value pair.")
			}
			key = fmt.Sprintf("%v%v", strings.ToUpper(key[0:1]), key[1:])

			// set flag
			lcc.SetInjected(key)
		} else {
			if debug {
				fmt.Printf(">>>> param: %10v = %v\n", key, param)
			}
			if key == "" || param == nil {
				// panic("value is nil")
				continue
			}

			// value, then set key to struct.
			v := reflect.ValueOf(param)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}

			if v.Kind() == reflect.String {
				// auto convert string into values.
				data[key] = []string{param.(string)}
			} else {
				// assign values into field.
				// TODO performance, cache this.
				injectField(lcc.V, key, param)
			}
		}
	}

	if len(data) > 0 {
		utils.SchemaDecoder.Decode(lcc.Proton, data)
	}
	return lcc
}

func injectField(target reflect.Value, fieldName string, value interface{}) {
	// target
	t := target
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// field
	// fmt.Printf(")))))))))))) : %v\n", fieldName)
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

	// fmt.Printf(">>>>>>>>>>>>>>>>> Inject %v::%v t:%v  <--  %v\n",
	// 	t, fieldName, field.Kind(), reflect.TypeOf(value),
	// )
	t.FieldByName(fieldName).Set(v)
}

func (lcc *LifeCircleControl) SetInjected(fields ...string) {
	method := lcc.V.MethodByName("SetInjected")
	if method.IsValid() {
		for _, f := range fields {
			method.Call([]reflect.Value{reflect.ValueOf(f), reflect.ValueOf(true)})
		}
	}
}
