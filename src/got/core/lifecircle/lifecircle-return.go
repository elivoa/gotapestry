/*
   Time-stamp: <[lifecircle-return.go] Elivoa @ Saturday, 2013-07-27 13:32:00>
*/
package lifecircle

import (
	"bytes"
	"errors"
	"fmt"
	"got/core"
	"got/templates"
	"net/http"
	"reflect"
	"strings"
)

/* ________________________________________________________________________________
   Return
*/

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
func (lcc *LifeCircleControl) Return(returns ...reflect.Value) bool {
	if returns == nil || len(returns) == 0 {
		return false // no return, continue
	}

	debuglog("... - [route.Return] handle return  '%v'", returns)
	if len(returns) > 0 {
		// determin kind
		returnValue := returns[0]
		kind := returnValue.Kind()
		if kind == reflect.Interface {
			returnValue = reflect.ValueOf(returnValue.Interface())
			kind = returnValue.Kind()
			// if is interface, reassign kind and value and continue
		}

		debuglog("... - [route.return] return type is '%v'", kind)

		// the first type is string. process string return.
		if kind == reflect.String {
			stringValue := returnValue.Interface().(string)
			lcc.ResultType = stringValue
			switch strings.ToLower(stringValue) {
			case "template":
				tname, err := extractString(1, returns...)
				if err != nil {
					lcc.Err = err
					return true
				}
				debuglog("-900- [route:return] parse template '%v'", tname)
				lcc.return_template(tname)

			case "text":
				debuglog("-902- [route:return] return plain string")
				lcc.return_text("plain/text", returns...)

			case "json":
				debuglog("-902- [route:return] return plain string")
				lcc.return_text("text/json", returns...)

			case "redirect":
				url, err := extractString(1, returns...)
				if err != nil {
					lcc.Err = err
					return true
				}
				debuglog("-904- [route:return] redirect to '%v'", url)
				http.Redirect(lcc.W, lcc.R, url, http.StatusFound)
			default:
				debuglog("[Warrning] return type %v not found!", stringValue)
				panic(fmt.Sprintf("[Warrning] return type %v not found!", stringValue))
			}
			return true
		} else if kind == reflect.Ptr {

			// may be the first return is an error.
			returnError, ok := returnValue.Interface().(error)
			if ok { // is an error
				// TODO create error page here.
				lcc.Err = returnError
				return true
			}
		}

	}
	return false
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

// support second parameter type: string, []byte
func (lcc *LifeCircleControl) return_text(contentType string, data ...reflect.Value) {
	// now we only return 1 result.
	// lcc.ResultType = data[1].String()
	v := data[1]
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	lcc.W.Header().Add("content-type", contentType)

	// write text
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		// only support byte slice
		lcc.W.Write(v.Bytes())
	case reflect.String:
		lcc.W.Write([]byte(v.String()))
	}
}

func (lcc *LifeCircleControl) return_template(templateName string) error {
	if templateName == "" {
		lcc.Err = errors.New(fmt.Sprintf(
			"%v %v must has an associated template or other return.",
			lcc.Kind, lcc.Name,
		))
		return lcc.Err
	}

	debuglog("-980- Render Tempalte %v.", templateName)

	var err error
	if lcc.Kind == core.COMPONENT {
		var buffer bytes.Buffer
		err = templates.RenderTemplate(&buffer, templateName, lcc.Proton)
		lcc.String = buffer.String()
	} else {
		err = templates.RenderTemplate(lcc.W, templateName, lcc.Proton)
	}

	if err != nil {
		lcc.Err = err
		return err
	}
	return nil
}
