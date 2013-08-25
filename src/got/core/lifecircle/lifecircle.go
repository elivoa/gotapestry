/*
   Time-stamp: <[lifecircle.go] Elivoa @ Saturday, 2013-08-24 14:50:03>
*/

package lifecircle

import (
	"bytes"
	"fmt"
	"github.com/gorilla/context"
	"got/core"
	"got/register"
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

// LifecircleControl controls flow of a Request.
type LifeCircleControl struct {
	// basic services
	w http.ResponseWriter
	r *http.Request

	// request related
	pageUrl   string // matched page url, Activate parameter part
	eventName string // an event call on page, not page render

	// lifes
	page    *Life // The Page Life
	current *Life // The Current Life

	// returns           || type:[template|redirect]
	rendering bool     // set to false to stop render.
	returns   *Returns //
	Err       error    // error if something error. TODO change to MultiError.

	// ResultType string // returns manually. if empty, find default tempalte
	// String     string // component html
}

// Life is a Page, Component, or others in the page render lifecircle.
type Life struct {
	// targets: are these too many?
	proton   core.Protoner // Pager or Componenter value
	rootType reflect.Type  // proton's root type
	v        reflect.Value // proton's reflect.Value
	kind     core.Kind     // enum: page|component
	name     string        // page name or component name

	registry *register.ProtonSegment // Is this really useful

	// tree structure. TODO need children?
	control   *LifeCircleControl
	container *Life
	embed     []*Life // not used

	// results
	out bytes.Buffer

	// no use?
	Path string // ???
}

// newControl create a new LifeCircleConstrol.
func newControl(w http.ResponseWriter, r *http.Request) *LifeCircleControl {
	lcc := &LifeCircleControl{w: w, r: r}
	return lcc
}

// createPage set the root page to lcc(LifeCircleControl).
func (lcc *LifeCircleControl) createPage(seed core.Protoner) *Life {
	life := newLife(seed)
	life.control = lcc
	lcc.page = life
	lcc.current = life
	debuglog("-710- [flow] New LifeCircleControl: %v[%v].", life.kind, life.name)
	return life
}

// appendComponent appends an embed component to lcc and chain it.
func (lcc *LifeCircleControl) appendComponent(seed core.Protoner) *Life {
	if seed.Kind() == core.PAGE {
		panic("Can't embed a Page!")
	}
	life := newLife(seed)
	life.control = lcc
	life.container = lcc.current
	lcc.current = life
	return lcc.current
}

// create new proton value from the seed.
func newLife(seed core.Protoner) *Life {
	life := &Life{}
	life.v = newProtonInstance(seed)
	life.proton = life.v.Interface().(core.Protoner)
	life.name = fmt.Sprint(reflect.TypeOf(life.proton).Elem()) // remove dependence of fmt
	life.rootType = utils.GetRootType(seed)
	life.kind = life.proton.Kind()
	return life
}

// ---- utils ---------------------------------------------------------------------

func (lcc *LifeCircleControl) SetToRequest(key interface{}, value interface{}) {
	context.Set(lcc.r, key, value)
}

func (lcc *LifeCircleControl) GetFromRequest(key interface{}) interface{} {
	return context.Get(lcc.r, key)
}

// ---- Accessors -----------------------------------------------------------------
func (lcc *LifeCircleControl) SetPageUrl(pageUrl string) *LifeCircleControl {
	lcc.pageUrl = pageUrl
	return lcc
}

func (lcc *LifeCircleControl) SetEventName(event string) *LifeCircleControl {
	lcc.eventName = event
	return lcc
}

// ---- Life ----------------------------------------------------------------------

func (l *Life) SetRegistry(registry *register.ProtonSegment) {
	l.registry = registry
}

// Call Events, and other events.
func (l *Life) call(names ...string) []reflect.Value {
	// fmt.Println("  ----  >> try call: ", names)
	// execute the first available method
	for _, name := range names {
		method := l.v.MethodByName(name)
		if method.IsValid() {
			debuglog("-730- [flow] Call Event: %v::%v().", l.name, name)
			returns := method.Call(emptyParameters)
			return returns
		}
	}
	return nil
}

// ********************************************************************************
// ********************************************************************************
// ********************************************************************************

// ________________________________________________________________________________

// POST Flow,
//   Events:
//     OnSubmit    - Form submitted, called before inject form values.
//     OnValidate  - Validate form. called after form value injected.
//                   If returns false, render the current page, with errors
//     OnSuccess   - Called if OnValidate returns true.
//
// TODO post to components.
func (lcc *LifeCircleControl) PostFlow() (returns *Returns) {
	// add ParseForm to fix bugs in go1.1.1
	err := lcc.r.ParseForm()
	if err != nil {
		lcc.Err = err
		return nil
	}

	// TODO use another method to retrive FormName. t:id
	formName := lcc.r.PostFormValue("t:id")
	if formName != "" {
		formName = fmt.Sprintf("From%v", formName)
	}
	fmt.Println("********************************************************************************")
	fmt.Println(formName)

	// call OnSubmit() method
	onSubmitEventName := fmt.Sprintf("%v%v", "OnSubmit", formName)
	returns = eventReturn(lcc.page.call(onSubmitEventName))
	if returns.breakReturn() {
		return
	}
	if returns.returnsFalse() {
		return lcc.refreshThisPage()
	}

	// inject form values
	lcc.InjectFormValues()

	// call OnValidate() method
	onValidateEventName := fmt.Sprintf("%v%v", "OnValidate", formName)
	returns = eventReturn(lcc.page.call(onValidateEventName))
	if returns.breakReturn() {
		return
	}
	if returns.returnsFalse() {
		return lcc.refreshThisPage()
	}

	// call success method
	// call OnSuccess() method
	onSuccessEventName := fmt.Sprintf("%v%v", "OnSuccess", formName)
	returns = eventReturn(lcc.page.call(onSuccessEventName))
	if returns.breakReturn() {
		return
	}
	fmt.Println("***************************************************************************")
	fmt.Println(onSuccessEventName)

	// something else, validation...
	// post flows stopd here.
	return lcc.refreshThisPage()
}

// ________________________________________________________________________________
// Call Events, and other events.
//
// func (lcc *LifeCircleControl) CallEvent(name string) bool {
// 	method := lcc.V.MethodByName(name)
// 	if method.IsValid() {
// 		debuglog("-730- [flow] Call Event: %v::%v().", lcc.Name, name)
// 		returns := method.Call(emptyParameters)
// 		return lcc.Return(returns...)
// 	} else {
// 		// debuglog("    - Event Not Found: %v", name)
// 	}
// 	return false
// }

// Call Events, with parameters. only used by Activate for now.
// TODO performance
func (lcc *LifeCircleControl) CallEventWithURLParameters(name string) bool {
	return lcc._callEventWithURLParameters(name, lcc.page.v)
}

func (lcc *LifeCircleControl) _callEventWithURLParameters(name string, base reflect.Value) bool {
	fmt.Println("______________________________________________________________")
	// fmt.Println(lcc.R.URL.Path)

	url := lcc.r.URL.Path
	if !strings.HasPrefix(url, lcc.pageUrl) {
		panic(fmt.Sprintf("%v should has prefix %v", url, lcc.pageUrl))
	}

	// parepare parameters, TODO extract method.
	paramsString := url[len(lcc.pageUrl)+1:]
	if lcc.eventName != "" {
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

		debuglog("-730- [flow] Call Event: %v::%v%v().", lcc.page.name, name, strParams)
		lcc.returns = eventReturn(method.Call(parameters))
		lcc.handleBreakReturn()
		return true
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
