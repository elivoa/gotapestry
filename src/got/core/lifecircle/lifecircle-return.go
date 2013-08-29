/*
   Time-stamp: <[lifecircle-return.go] Elivoa @ Tuesday, 2013-08-27 13:13:38>
*/
package lifecircle

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

/* ________________________________________________________________________________
   Return
*/

// results
type Returns struct {
	returnType string // bool, redirect, text, json
	value      interface{}
}

func (r *Returns) breakReturn() bool {
	// returns true or nothing
	if r.returnType == "bool" && r.value == true {
		return false
	}
	if r.returnType == "template" {
		return false
	}
	return true
}

func (r *Returns) returnsTrue() bool  { return r.returnType == "bool" && r.value == true }
func (r *Returns) returnsFalse() bool { return r.returnType == "bool" && r.value == false }

func (lcc *LifeCircleControl) refreshThisPage() *Returns {
	return &Returns{
		returnType: "redirect",
		value:      lcc.r.URL.Path,
	}
}

/*
   Return analysis the return values and redirect to the write place.

   @return: true if is redirect and should return.
            false if continue the flow.

   Now can processing: string, string interface...

   Process the return values:
   If the first value is an error. Redirect to Error page.
   If the only one value is string, confirm this a url and rediredt.
   If has at least 2 values and first parameter is the following string.
      template - template path
      json     - json string / json bytes
      text     - plan/text
      redirect - redirect to url
   @Author elivoa@gmail.com
*/
// event call returns
func eventReturn(returns []reflect.Value) *Returns {

	// debuglog("... - [route.Return] handle return  '%v'", returns)

	// returns nothing equals return true
	if returns == nil || len(returns) == 0 {
		return returnTrue()
	}

	first := returns[0]
	kind := first.Kind()
	// dereference interface
	if first.Kind() == reflect.Interface {
		first = reflect.ValueOf(first.Interface())
		kind = first.Kind()
	}

	switch kind {
	case reflect.Bool:
		return &Returns{"bool", returns[0].Bool()}
	case reflect.String:
		// now only support (string, string) pair
		if len(returns) <= 1 {
			panic("return string must have the second return value.")
		}
		ft := strings.ToLower(first.String())
		switch ft {
		case "text", "json":
			return &Returns{ft, returns[1].Interface()}
		case "redirect":
			url, err := extractString(1, returns...)
			if err != nil {
				panic(err.Error())
			}
			return &Returns{"redirect", url}
		// case "template":
		default:
			debuglog("[Warrning] return type %v not found!", first.String())
			panic(fmt.Sprintf("[Warrning] return type %v not found!", first.String()))
		}

	case reflect.Invalid: // invalid means return nil
		return returnTrue()

	case reflect.Ptr:
		// may be the first return is an error.
		err, ok := first.Interface().(error)
		if ok { // is an error
			panic(err.Error())
		}
	}
	return nil
}

func returnTrue() *Returns  { return &Returns{"bool", true} }
func returnFalse() *Returns { return &Returns{"bool", false} }

func (lcc *LifeCircleControl) handleBreakReturn() {
	r := lcc.returns
	if r == nil {
		return
	}
	switch r.returnType {
	case "text":
		lcc.return_text("plain/text", r.value)
	case "json":
		lcc.return_text("text/json", r.value)
	case "redirect":
		// debuglog("-904- [route:return] redirect to '%v'", url)
		http.Redirect(lcc.w, lcc.r, r.value.(string), http.StatusFound)
	}
}

// support second parameter type: string, []byte
func (lcc *LifeCircleControl) return_text(contentType string, value interface{}) {
	// now we only return 1 result.
	lcc.w.Header().Add("content-type", contentType)
	lcc.w.Write([]byte(fmt.Sprint(value)))

	// TODO struct to json

	// // write text
	// switch v.Kind() {
	// case reflect.Slice, reflect.Array:
	// 	// only support byte slice
	// 	lcc.W.Write(v.Bytes())
	// case reflect.String:
	// 	lcc.W.Write([]byte(v.String()))
	// }
}

// TODO reflect utils.
func extractString(index int, data ...reflect.Value) (string, error) {
	if len(data) <= index {
		err := errors.New(fmt.Sprintf("Need the %v(st/nd) return value."))
		return emptyString, err
	}
	if data[index].Kind() == reflect.String {
		return data[index].Interface().(string), nil
	} else {
		err := errors.New(fmt.Sprintf("The %v(st/nd) return value must be string."))
		return emptyString, err
	}
}

// func eventReturn___________(returns []reflect.Value) *Returns {
// 	// ********************************************************************************
// 	if len(returns) > 0 {
// 		// determin kind
// 		returnValue := returns[0]
// 		kind := returnValue.Kind()
// 		if kind == reflect.Interface {
// 			returnValue = reflect.ValueOf(returnValue.Interface())
// 			kind = returnValue.Kind()
// 			// if is interface, reassign kind and value and continue
// 		}

// 		debuglog("... - [route.return] return type is '%v'", kind)

// 		// the first type is string. process string return.
// 		if kind == reflect.String {
// 			stringValue := returnValue.Interface().(string)
// 			lcc.ResultType = stringValue
// 			switch strings.ToLower(stringValue) {
// 			// case "template":
// 			// 	tname, err := extractString(1, returns...)
// 			// 	if err != nil {
// 			// 		lcc.Err = err
// 			// 		return true
// 			// 	}
// 			// 	debuglog("-900- [route:return] parse template '%v'", tname)
// 			// 	lcc.TemplateName = tname

// 			case "text":
// 				debuglog("-902- [route:return] return plain string")
// 				lcc.return_text("plain/text", returns...)

// 			case "json":
// 				debuglog("-902- [route:return] return plain string")
// 				lcc.return_text("text/json", returns...)

// 			case "redirect":
// 				url, err := extractString(1, returns...)
// 				if err != nil {
// 					lcc.Err = err
// 					return true
// 				}
// 				debuglog("-904- [route:return] redirect to '%v'", url)
// 				http.Redirect(lcc.W, lcc.R, url, http.StatusFound)
// 			default:
// 				debuglog("[Warrning] return type %v not found!", stringValue)
// 				panic(fmt.Sprintf("[Warrning] return type %v not found!", stringValue))
// 			}
// 			return true
// 		} else if kind == reflect.Ptr {

// 			// may be the first return is an error.
// 			returnError, ok := returnValue.Interface().(error)
// 			if ok { // is an error
// 				// TODO create error page here.
// 				lcc.Err = returnError
// 				return true
// 			}
// 		}

// 	}
// 	return false

// }

// ********************************************************************************
// ********************************************************************************
// ********************************************************************************
// ********************************************************************************
// ********************************************************************************
// ********************************************************************************

// func (lcc *LifeCircleControl) Return(returns ...reflect.Value) bool {
// }

// // if no return value or true, returns true
// // panic if return type is not bool
// func eventReturnBool(returns []reflect.Value) bool {
// 	if returns == nil || len(returns) == 0 {
// 		return true
// 	}
// 	if len(returns) > 1 || returns[0].Kind() != reflect.Bool {
// 		panic("return value must be bool or no return value!")
// 	}
// 	return returns[0].Bool()
// }
