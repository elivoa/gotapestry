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
	"path"
	"reflect"
	"strings"
)

/* ________________________________________________________________________________
   TODO:
   . Callback functions
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
	Kind      core.Kind     // enum: page|component
	V         reflect.Value // Value of page
	Name      string        // page name or component name
	Path      string        // ???
	PageUrl   string        // matched page url, Activate parameter part
	EventName string        // an event call on page, not page render

	// result type:[template|redirect]
	ResultType string // returns manually. if empty, find default tempalte
	String     string // component html
	Err        error  // error if something error
}

// --------------------------------------------------------------------------------
// Create new page value with the same type as proton (page or component).
// by calling New method or use reflect.
func NewPageFlow(w http.ResponseWriter, r *http.Request, page core.Pager) *LifeCircleControl {
	// maintaince structCache
	if si := scache.GetPageX(reflect.TypeOf(page)); si == nil {
		panic("Can't parse page!")
	}

	return newLifeCircleControl(w, r, core.PAGE, page)
}

// Create a new Component Flow.
// param:
//   container - real container object.
//   component - current component base object.
//   params - parameters in the component grammar.
//
// Note: I maintain StructCache here in the flow create func. This occured only when
//       page or component are rendered. Directly post to a page can not invoke structcache init.
//
func NewComponentFlow(container core.Protoner, component core.Componenter,
	params []interface{}) *LifeCircleControl {

	debuglog("----- [Create Component flowcontroller] ------------------------%v",
		"----------------------------------------")
	debug.Log("- C - [Component Container] Type: %v, ComponentType:%v,\n",
		reflect.TypeOf(container), reflect.TypeOf(component))

	// Store type in StructCache, Store instance in ProtonObject.
	// Warrning: What if use component in page/component but it's not initialized?
	// Tid= xxx in template must the same with fieldname in .go file.
	//

	// 1. cache in StructInfoCache. (application scope)
	si := scache.GetCreate(reflect.TypeOf(container), container.Kind())
	if si == nil {
		panic(fmt.Sprintf("StructInfo for %v can't be null!", reflect.TypeOf(container)))
	}
	t, _ := utils.RemovePointer(reflect.TypeOf(component), false)
	tid, _ := determinComponentTid(params, t)
	si.CacheEmbedProton(t, tid, component.Kind())

	// 2. store in proton's embed field. (request scope)
	proton, ok := container.Embed(tid)
	var lcc *LifeCircleControl
	if !ok {
		// The first create new component object.
		lcc = newLifeCircleControl(container.ResponseWriter(), container.Request(),
			core.COMPONENT, component)
		container.SetEmbed(tid, lcc.Proton)
	} else {
		// If proton in a loop, we use the same proton instance.
		lcc = &LifeCircleControl{
			W:      container.ResponseWriter(),
			R:      container.Request(),
			Kind:   core.COMPONENT,
			Proton: proton,
			V:      reflect.ValueOf(proton),
			Name:   fmt.Sprint(reflect.TypeOf(proton).Elem()),
		}
		proton.IncEmbed() // increase loop index. Used by ClientId()
	}

	lcc.InjectComponentParameters(params) // inject component parameters
	return lcc
}

// return name, is setManually; t must not be ptr.
func determinComponentTid(params []interface{}, t reflect.Type) (tid string, setManually bool) {
	for idx, p := range params {
		if idx%2 == 0 && strings.ToLower(p.(string)) == "tid" {
			tid = params[idx+1].(string)
		}
	}
	if tid == "" {
		setManually = true
		tid = path.Ext(t.String())[1:]
	}
	return
}

func newLifeCircleControl(w http.ResponseWriter, r *http.Request,
	kind core.Kind, proton core.Protoner) *LifeCircleControl {

	lcc := &LifeCircleControl{W: w, R: r, Kind: kind}
	lcc.V = newProtonInstance(proton)
	lcc.Proton = lcc.V.Interface().(core.Protoner)
	lcc.Name = fmt.Sprint(reflect.TypeOf(lcc.Proton).Elem())

	debuglog("-710- [flow] New LifeCircleControl: %v[%v].", lcc.Kind, lcc.Name)

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

	// Only Page has Activate event.
	if lcc.Kind == core.PAGE {
		// Life Circle: Activate (TODO extract function)
		// if ret := lcc.CallEventWithURLParameters("Activate"); ret {
		if ret := lcc.CallEvent("Activate"); ret {
			return lcc
		}
	}

	// On form submit, go to another flow.
	// TODO：只有page才POST，修复component中提交form的bug。设计一个方案。
	if lcc.R.Method == "POST" && lcc.Kind == core.PAGE {
		return lcc.PostFlow()
	}

	// LifeCircle
	for _, eventName := range flowEvents {
		if ret := lcc.CallEvent(eventName); ret {
			return lcc
		}
	}
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
// event call page flow
// TODO support component event call
func (lcc *LifeCircleControl) EventCall(event string) *LifeCircleControl {
	// 1. Inject values into root page
	lcc.InjectValue()

	// 2. Call Activate method on root page
	if lcc.Kind == core.PAGE {
		if ret := lcc.CallEvent("Activate"); ret {
			return lcc
		}
	}

	proton := lcc.Proton // proton is father node.
	eventPaths := strings.Split(event, ".")
	for idx, piece := range eventPaths {
		if idx < len(eventPaths)-1 { // call path
			// 1. get from proton cache; !!!This can't be happened!!!
			c, ok := proton.Embed(piece)
			if ok && c != nil {
				proton = c
				continue
			}

			// 2. Is Cached in StructInfo
			si := scache.GetCreate(reflect.TypeOf(proton), proton.Kind()) // root page
			if si == nil {
				panic(fmt.Sprintf("StructInfo for %v can't be null!", reflect.TypeOf(proton)))
			}

			var newProton core.Protoner
			fi := si.FieldInfo(piece)
			if fi != nil {
				newInstance := newInstance(fi.Type)
				newProton = newInstance.Interface().(core.Protoner)
			} else {
				// If not cached fieldInfo, create FieldInfo
				containerType, _ := utils.RemovePointer(reflect.TypeOf(proton), false)
				field, ok := containerType.FieldByName(piece)
				if !ok {
					panic(fmt.Sprintf("Can't get field in path: %v", piece))
				}
				newInstance := newInstance(field.Type)                   // create new instance
				newProton = newInstance.Interface().(core.Protoner)      //
				si.CacheEmbedProton(field.Type, piece, newProton.Kind()) // cache
			}
			proton.SetEmbed(piece, newProton) // store newInstance into proton
			proton = newProton                // next round

		} else { // last node
			// Call event. TODO add parameters.
			fmt.Println("\n----------    EVENT CALL    ----------------")
			fmt.Printf("Call event [%v] with parameters.\n", event)
			if ret := lcc._callEventWithURLParameters("On"+piece, reflect.ValueOf(proton)); ret {
				return lcc
			}
		}
	}
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
		// TODO redirect to error page.
		lcc.Err = err
		return err
	}
	return nil
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
