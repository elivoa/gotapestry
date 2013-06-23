package register

import (
	"fmt"
	"got/core"
	"got/templates"
	"got/utils"
	"log"
	"reflect"
	"strings"
)

/*_______________________________________________________________________________
  Component Registry
*/

// value
var ccache *SimpleComponentCache = &SimpleComponentCache{map[string]core.IComponent{}}

// only support components under root folder.
type SimpleComponentCache struct {
	Components map[string]core.IComponent
}

// add to cache
func (s *SimpleComponentCache) Add(c core.IComponent) string {
	// 1. add to cache
	name := reflect.TypeOf(c).Elem().Name()
	s.Components[name] = c

	// 2. register template
	templates.Add(fmt.Sprintf("components/%v", name))

	// 3. register template-func
	templates.RegisterComponent(strings.ToLower(name), componentLifeCircle(name))
	return name
}

// register method
func RegisterComponent(c core.IComponent) {
	log.Printf("[building] Register component: %v\n", reflect.TypeOf(c))

	name := ccache.Add(c)

	log.Printf("-------------------------------------------------------")
	log.Printf("key is : %v", name)
}

// not used
func newComponentValue(base interface{}) reflect.Value {
	baseValue := reflect.ValueOf(base)
	if baseValue.Kind() == reflect.Ptr {
		// pass?
	}
	newValue := reflect.New(reflect.TypeOf(base).Elem())

	fmt.Println("------------------------")
	fmt.Println(reflect.TypeOf(base))
	fmt.Println(reflect.TypeOf(newValue))
	fmt.Println(newValue)
	return newValue
}

func Unmarshal(target interface{}, params ...interface{}) error {
	data := make(map[string][]string)
	var key string
	for i, param := range params {
		if i%2 == 0 {
			key = param.(string)
			key = fmt.Sprintf("%v%v", strings.ToUpper(key[0:1]), key[1:])
		} else {
			// value, then set key to struct.
			v := reflect.ValueOf(target)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			if v.Kind() == reflect.String {
				// auto convert string into values.
				data[key] = []string{param.(string)}
			} else {
				// assign values into field.
				// TODO performance, cache this.
				field := v.FieldByName(key)
				field.Set(reflect.ValueOf(param))
			}
		}
	}
	if len(data) > 0 {
		utils.SchemaDecoder.Decode(target, data)
	}
	return nil
}
