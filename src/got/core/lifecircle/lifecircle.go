package lifecircle

import (
	"bytes"
	"errors"
	"fmt"
	"got/core"
	"got/debug"
	"got/templates"
	"got/utils"
	"log"
	"net/http"
	"reflect"
	"strings"
)

/* ________________________________________________________________________________
   TODO:
   . Callback functions
*/
const (
	// TemplateKey = "template-key"
	PageKey  = "page-key"
	debugLog = true

	PAGE      = "page"
	COMPONENT = "component"
)

func debuglog(format string, params ...interface{}) {
	if debugLog {
		log.Printf(format, params...)
	}
}

var (
	emptyParameters = []reflect.Value{}
	emptyString     = ""
)

/* ______________________________
   Control Struct
*/
type LifeCircleControl struct {
	// basic service
	W http.ResponseWriter
	R *http.Request

	// target
	Proton    core.IProton  // IPage or IComponent value
	V         reflect.Value // Value of page
	Name      string        // page name or component name
	Path      string        // ???
	Kind      string        // enum: page|component
	PageUrl   string        // matched page url, Activate parameter part
	EventName string        // an event call on page, not page render

	// result type:[template|redirect]
	ResultType string // returns manually. if empty, find default tempalte
	String     string // component html
	Err        error  // error if something error
}

/* ______________________________
   Create
*/
func newLifeCircleControl(
	w http.ResponseWriter, r *http.Request,
	kind string, proton core.IProton) *LifeCircleControl {

	lcc := &LifeCircleControl{W: w, R: r, Kind: kind}
	baseValue := reflect.ValueOf(proton)
	method := baseValue.MethodByName("New")
	if method.IsValid() {
		// create by calling New
		returns := method.Call(emptyParameters)
		lcc.V = returns[0]
	} else {
		// create an empty page value
		lcc.V = reflect.New(reflect.TypeOf(proton).Elem())
	}
	lcc.Proton = lcc.V.Interface().(core.IProton)
	lcc.Name = fmt.Sprint(reflect.TypeOf(lcc.Proton).Elem())

	debuglog("-710- [flow] New LifeCircleControl: %v[%v].",
		lcc.Kind, lcc.Name)

	return lcc
}

// Create new page value with the same type as proton (page or component).
// by calling New method or use reflect.
func NewPageFlow(
	w http.ResponseWriter, r *http.Request, proton core.IPage) *LifeCircleControl {

	lcc := newLifeCircleControl(w, r, PAGE, proton)
	return lcc
}

// Method to create a new Component value.
func NewComponentFlow(
	container core.IProton, component core.IComponent,
	params []interface{}) *LifeCircleControl {

	fmt.Println(".................... [Cretae Component Flow]" +
		" ---------------------------------------------------------------")
	debug.Log("- C - [Component Container] Type: %v, ComponentType:%v,\n",
		reflect.TypeOf(container), reflect.TypeOf(component))

	lcc := newLifeCircleControl(container.ResponseWriter(), container.Request(),
		COMPONENT, component)
	// inject component parameters

	lcc.InjectComponentParameters(params)
	return lcc
}

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

/*
   Call LifeCircles
   Page lifecircle
*/
var flowEvents = []string{
	"Setup",
	"SetupRender", // SetupRender is deprecated. use Setup instead.
	"BeginRender",
	"AfterRender",
}

// normal page flow
func (lcc *LifeCircleControl) Flow() *LifeCircleControl {
	lcc.InjectValue()

	// ____________________________________________________________________
	// if false { // disable Activte, this has difficult to implement.
	if lcc.Kind == "page" {
		// Life Circle: Activate (TODO extract function)
		// if ret := lcc.CallEventWithURLParameters("Activate"); ret {
		if ret := lcc.CallEvent("Activate"); ret {
			return lcc
		}
	}
	// }

	// On form submit, execute above to create page value.
	// ignore the following and call OnSuccess event.
	// fmt.Printf("r.Method is %v\n", r.Method)
	// TODO:
	//    Support OnSuccess Method.
	// TODO：只有page才POST，修复component中提交form的bug。设计一个方案。
	if lcc.R.Method == "POST" && lcc.Kind == "page" {
		return lcc.PostFlow()
	}

	// LifeCircle
	for _, eventName := range flowEvents {
		if ret := lcc.CallEvent(eventName); ret {
			return lcc
		}
	}

	return lcc // for chain
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

// event call page flow
// TODO support component event call
func (lcc *LifeCircleControl) EventCall(event string) *LifeCircleControl {
	lcc.InjectValue()
	if lcc.Kind == "page" {
		// Life Circle: Activate (TODO extract function)
		// if ret := lcc.CallEventWithURLParameters("Activate"); ret {
		if ret := lcc.CallEvent("Activate"); ret {
			return lcc
		}
	}

	// call event. TODO add parameters.
	fmt.Println("\n\n\n$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	fmt.Println("Call event %v with parameters.", event)
	if ret := lcc.CallEventWithURLParameters("On" + event); ret {
		return lcc
	}
	return lcc // for chain
}

// Call Events, and other events.
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

	fmt.Printf("-0- params is %v\n", paramsString)
	strParams := strings.Split(paramsString, "/")
	for _, strParam := range strParams {
		// TODO inject values.
		fmt.Printf("-1- param is: %v\n", strParam)
	}

	//reflect.TypeOf(method).NumIn
	method := lcc.V.MethodByName(name)
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
	if lcc.Kind == "component" {
		var buffer bytes.Buffer
		err = templates.RenderTemplate(&buffer, templateName, lcc.Proton)
		lcc.String = buffer.String()
	} else {
		err = templates.RenderTemplate(lcc.W, templateName, lcc.Proton)
	}

	if err != nil {
		// TODO redirect to error page.
		lcc.Err = err
		return err
	}
	return nil
}
