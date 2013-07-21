package lifecircle

import (
	"got/core"
	"reflect"
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

func SetInjected(v reflect.Value, fields ...string) {
	method := v.MethodByName("SetInjected")
	if method.IsValid() {
		for _, f := range fields {
			method.Call([]reflect.Value{reflect.ValueOf(f), reflect.ValueOf(true)})
		}
	}
}
