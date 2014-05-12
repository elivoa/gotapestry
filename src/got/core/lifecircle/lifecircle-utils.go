package lifecircle

import (
	"fmt"
	"got/core"
	"got/debug"
	"reflect"
	"strings"
)

// create new proton instance with base protoner
func newProtonInstance(proton core.Protoner) reflect.Value {
	baseValue := reflect.ValueOf(proton)

	// try to create new value of proton
	method := baseValue.MethodByName("New")
	if method.IsValid() {
		returns := method.Call(emptyParameters)
		return returns[0]
	} else {
		// return reflect.New(reflect.TypeOf(proton).Elem())
		return newInstance(reflect.TypeOf(proton))
	}
}

// create new instance by type.
func newInstance(rt reflect.Type) reflect.Value {
	t := rt
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return reflect.New(t)
}

// return parameters
func _extractParameters(url string, pageUrl string, eventName string) []string {
	return nil
}

// maybe not used.
// TODO performance, maybe i can directly call method. don't use reflect.
func SetInjected(v reflect.Value, fields ...string) {
	method := v.MethodByName("SetInjected")
	if method.IsValid() {
		for _, f := range fields {
			method.Call([]reflect.Value{reflect.ValueOf(f), reflect.ValueOf(true)})
		}
	}
}

// TODO Coercion: make an interface to auto translate this into it.
func analysisTranslateSuffix(t reflect.Type) string {
	switch t.String() {
	case "*gxl.Int":
		return ".Int"
	}
	return ""
}

// analysis request url, split parameters from page.
// param:
//   url - full url
//   pageUrl - url represent a page
func extractPathParameters(url string, pageUrl string, eventName string) []string {
	// validate
	if !strings.HasPrefix(url, pageUrl) {
		panic(fmt.Sprintf("%v should has prefix %v", url, pageUrl))
	}

	// parepare parameters
	paramsString := url[len(pageUrl):]
	if eventName != "" {
		index := strings.Index(paramsString, "/")
		if index > 0 {
			paramsString = paramsString[index:]
		}
	}
	var pathParams []string
	if len(paramsString) > 0 {
		if strings.HasPrefix(paramsString, "/") {
			paramsString = paramsString[1:]
		}
		pathParams = strings.Split(paramsString, "/")
	}
	debug.Log("-   - [injection] URL:%v, parameters:%v", url, pathParams)
	return pathParams
}

