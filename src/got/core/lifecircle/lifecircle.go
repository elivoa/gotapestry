/*
   Time-stamp: <[lifecircle.go] Elivoa @ Saturday, 2013-07-27 10:51:22>
*/

package lifecircle

import (
	"fmt"
	"got/core"
	"got/utils"
	"log"
	"net/http"
	"reflect"
	"strings"
)

/* ________________________________________________________________________________

Injection Order:
  Page Render:       new -> path -> url ->
  Event Call:        new -> page's path -> page's url -> event func's parameter -> call
  Component Render:  new -> page's path -> page's url -> component parameter
  Form Submit:       new -> page's path -> page's url -> form values

TODO:
 . Callback functions
 . Cache Method to inject, borrowed form gorilla/schema


*/

// ______________________________
// Control Struct
//
type LifeCircleControl struct {
	// basic service
	W http.ResponseWriter
	R *http.Request

	// target
	Proton    core.Protoner // Pager or Componenter value
	RootType  reflect.Type  // proton's root type
	V         reflect.Value // proton's reflect.Value
	Kind      core.Kind     // enum: page|component
	Name      string        // page name or component name
	Path      string        // ???
	PageUrl   string        // matched page url, Activate parameter part
	EventName string        // an event call on page, not page render

	// result type:[template|redirect]
	ResultType string // returns manually. if empty, find default tempalte
	String     string // component html
	Err        error  // error if something error. TODO change to MultiError.
}

// TODO kill this
func newLifeCircleControl(w http.ResponseWriter, r *http.Request,
	kind core.Kind, proton core.Protoner) *LifeCircleControl {

	lcc := &LifeCircleControl{W: w, R: r, Kind: kind}
	lcc.V = newProtonInstance(proton)
	lcc.Proton = lcc.V.Interface().(core.Protoner)
	lcc.Name = fmt.Sprint(reflect.TypeOf(lcc.Proton).Elem())
	lcc.RootType = utils.GetRootType(proton)

	debuglog("-710- [flow] New LifeCircleControl: %v[%v].", lcc.Kind, lcc.Name)

	return lcc
}

//
// Page lifecircle methods
//

// --------------------------------------------------------------------------------
// TODO Refactor NewPageFlow, NewComponentFlow

/* ______________________________
   Flow
*/
func (lcc *LifeCircleControl) SetPageUrl(pageUrl string) *LifeCircleControl {
	lcc.PageUrl = pageUrl
	return lcc
}

func (lcc *LifeCircleControl) SetEventName(event string) *LifeCircleControl {
	lcc.EventName = event
	return lcc
}

// ________________________________________________________________________________
// POST Flow,
//   Events:
//     OnSubmit    - Form submitted, called before inject form values.
//     OnValidate  - Validate form. called after form value injected.
//                   If returns false, render the current page, with errors
//     OnSuccess   - Called if OnValidate returns true.
//
func (lcc *LifeCircleControl) PostFlow() *LifeCircleControl {
	// add ParseForm to fix bugs in go1.1.1
	err := lcc.R.ParseForm()
	if err != nil {
		lcc.Err = err
		return lcc
	}

	// TODO use another method to retrive FormName. t:id
	formName := lcc.R.PostFormValue("t:id")
	if formName != "" {
		formName = fmt.Sprintf("From%v", formName)
	}
	fmt.Println("********************************************************************************")
	fmt.Println(formName)
	// call OnSubmit() method
	onSubmitEventName := fmt.Sprintf("%v%v", "OnSubmit", formName)
	if lcc.CallEvent(onSubmitEventName) {
		return lcc
	}

	// inject form values
	lcc.InjectFormValues()

	// call OnValidate() method
	onValidateEventName := fmt.Sprintf("%v%v", "OnValidate", formName)
	if lcc.CallEvent(onValidateEventName) {
		return lcc
	}

	// call success method
	// call OnSuccess() method
	onSuccessEventName := fmt.Sprintf("%v%v", "OnSuccess", formName)
	fmt.Println("********************************************************************************")
	fmt.Println(onSuccessEventName)

	if lcc.CallEvent(onSuccessEventName) {
		return lcc
	}

	// something else, validation...
	// post flows stopd here.
	return lcc
}

// ________________________________________________________________________________
// Call Events, and other events.
//
func (lcc *LifeCircleControl) CallEvent(name string) bool {
	method := lcc.V.MethodByName(name)
	if method.IsValid() {
		debuglog("-730- [flow] Call Event: %v::%v().", lcc.Name, name)
		returns := method.Call(emptyParameters)
		return lcc.Return(returns...)
	} else {
		// debuglog("    - Event Not Found: %v", name)
	}
	return false
}

// Call Events, with parameters. only used by Activate for now.
// TODO performance
func (lcc *LifeCircleControl) CallEventWithURLParameters(name string) bool {
	return lcc._callEventWithURLParameters(name, lcc.V)
}

func (lcc *LifeCircleControl) _callEventWithURLParameters(name string, base reflect.Value) bool {
	fmt.Println("______________________________________________________________")
	// fmt.Println(lcc.R.URL.Path)

	url := lcc.R.URL.Path
	if !strings.HasPrefix(url, lcc.PageUrl) {
		panic(fmt.Sprintf("%v should has prefix %v", url, lcc.PageUrl))
	}

	// parepare parameters, TODO extract method.
	paramsString := url[len(lcc.PageUrl)+1:]
	if lcc.EventName != "" {
		index := strings.Index(paramsString, "/")
		if index > 0 {
			paramsString = paramsString[index+1:]
		}
	}

	strParams := strings.Split(paramsString, "/")
	for idx, strParam := range strParams {
		// TODO inject values.
		fmt.Printf("-1- param #%d is: %v\n", idx, strParam)
	}

	//reflect.TypeOf(method).NumIn
	// method := lcc.V.MethodByName(name)
	method := base.MethodByName(name)
	if method.IsValid() {
		t := method.Type()
		fmt.Printf("-2- method is: %v\n", method)
		fmt.Printf("-2- typeof method is: %v\n", t)
		fmt.Printf("-2- param number is: %v\n", t.NumIn())

		numOfParameters := t.NumIn()
		parameters := make([]reflect.Value, numOfParameters)
		for i := 0; i < numOfParameters; i++ {
			pt := t.In(i)
			fmt.Printf("-3- param %v is: %v\n", i, pt)
			if len(strParams) > i {
				v, err := utils.Coercion(strParams[i], pt)
				if err != nil {
					panic(err.Error())
				}
				parameters[i] = v
			} else {
				parameters[i] = utils.CoercionNil(pt)
			}
		}

		debuglog("-730- [flow] Call Event: %v::%v%v().", lcc.Name, name, strParams)
		returns := method.Call(parameters)
		return lcc.Return(returns...)
	} else {
		// debuglog("    - Event Not Found: %v", name)
	}
	return false
}

// --------------------------------------------------------------------------------

const (
	debugLog = true
)

var (
	emptyParameters = []reflect.Value{}
	emptyString     = ""
)

func debuglog(format string, params ...interface{}) {
	if debugLog {
		log.Printf(format, params...)
	}
}
