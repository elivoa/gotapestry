package route

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"got/core"
	"got/core/lifecircle"
	"got/debug"
	"got/register"
	"got/templates"
	"html/template"
	"log"
	"net/http"
	"reflect"
	rd "runtime/debug"
	"strings"
)

var (
	emptyParameters = []reflect.Value{}
	debugLog        = true
)

// RouteHandler is responsible to handler all got request.
func RouteHandler(w http.ResponseWriter, r *http.Request) {
	url := "/" + mux.Vars(r)["url"]

	// 1. skip special resources. TODO Expand to config.
	if url == "/favicon.ico" {
		return
	}

	printAccessHeader(r)

	// --------  Error Handling  --------------------------------------------------------------
	defer func() {
		if err := recover(); err != nil {
			processPanic(err, r)
			// TODO Render error page 500/404 page.
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		}
		printAccessFooter(r)
	}()
	// --------  Routing...  --------------------------------------------------------------

	// 3. let's find the right pages.
	result := lookup(url)
	if result == nil {
		// TODO goto 404 page.
		// TODO all path-parameter parse error goes to 404, not 500.
	}
	debug.Log("-601- [RouteFind] %v", result.Segment)

	// TODO: Create New page object every request? howto share some page object? see tapestry5.
	lcc := lifecircle.NewPageFlow(w, r, result.Segment.Proton)
	lcc.SetPageUrl(result.PageUrl)
	lcc.SetEventName(result.EventName)

	if result.EventName == "" {
		// page render flow
		lcc.Flow()
		handleReturn(lcc, result.Segment)
	} else {
		// event call
		lcc.EventCall(result.EventName)

		// TODO wrong here. this is wrong. sudo refactor lifecircle-return.
		// if lcc not returned, return the current page.
		if lcc.Err != nil {
			panic(lcc.Err.Error())
		}
		// default return the current page.
		if result.Segment != nil && lcc.ResultType == "" {
			url := lcc.R.URL.Path
			http.Redirect(lcc.W, lcc.R, url, http.StatusFound)
		}
	}
}

func lookup(url string) *register.LookupResult {
	result, err := register.Pages.Lookup(url)
	if nil != err {
		panic(err.Error())
	}
	if result == nil || result.Segment == nil {
		panic(fmt.Sprintf("Error: seg.Proton is null. seg: %v", result.Segment))
	}
	if result.Segment.Proton == nil {
		// TODO redirect to 404 page.
		panic(fmt.Sprintf("~~~~ Page not found ~~~~"))
	}
	return result
}

// handle return
func handleReturn(lcc *lifecircle.LifeCircleControl, seg *register.ProtonSegment) {
	// no error, no templates return or redirect.
	if seg != nil && lcc.Err == nil && lcc.ResultType == "" {
		// find default tempalte to return
		identity, templatePath := seg.TemplatePath()
		// debug.Log("-755- [TemplateSelect] %v -> %v", key, tplPath)
		if _, err := templates.GotTemplateCache.Get(identity, templatePath); err != nil {
			lcc.Err = err
		} else {
			if err := templates.RenderGotTemplate(lcc.W, identity, lcc.Proton); err != nil {
				lcc.Err = err
			}
		}
	}
	// panic here if has errors.
	if lcc.Err != nil {
		panic(lcc.Err.Error())
	}
}

// --------------------------------------------------------------------------------

// handle components

// --------------------------------------------------------------------------------

// helper
func printAccessHeader(r *http.Request) {
	fmt.Println()
	fmt.Println(".......................................        " + r.URL.Path +
		"        ...............................................")
	// log.Printf(">>> access %v\n", r.URL.Path)
	// log.Printf("> w is %v\n", reflect.TypeOf(w))
	// log.Printf("> w is %v\n", reflect.TypeOf(req))
}

func printAccessFooter(r *http.Request) {
	//debug.Log("^ ^ ^ ^ ^ ^ ^ ^ PAGE RENDER END ^ ^ ^ ^ ^ ^ ^ ^ ^ ^")
	fmt.Println("-----------------------------^         PAGE RENDER END           " +
		"-----------------------------------")
	fmt.Println("................................................................." +
		"...................................")
}

func processPanic(err interface{}, r *http.Request) {
	log.Print("xxxxxxxx  PANIC  xxxxxxxxxxxxx", yibaix)
	log.Printf("x URL: %-109v x", r.URL.Path)
	log.Printf("x panic: %-109v x", err)
	log.Print("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", yibaix)
	fmt.Println("StackTrace >>")
	rd.PrintStack()
	fmt.Println()
}

var yibaix = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

// --------------------------------------------------------------------------------

func Page(f func(), pages ...core.Pager) int {
	return register.Page(f, pages...)
}

func Component(f func(), components ...core.Componenter) int {
	for _, c := range components {
		url := register.MakeUrl(f, c)
		// TODO has space to improve.
		selectors := register.Components.Add(url, c)
		for _, selector := range selectors {
			lowerKey := strings.ToLower(strings.Join(selector, "/"))
			templates.RegisterComponent(lowerKey, componentLifeCircle(lowerKey))
		}
	}
	return len(components)
}

/* ________________________________________________________________________________
   Execute Components
*/

/*
  Component Render Handler method
  Return: string or template.HTML
*/
func componentLifeCircle(name string) func(...interface{}) interface{} {

	return func(params ...interface{}) interface{} {

		log.Printf("-620- [flow] Render Component %v ....", name)

		// 1. find base component type
		result, err := register.Components.Lookup(name)
		if err != nil || result.Segment == nil {
			panic(fmt.Sprintf("Component %v not found!", name))
		}
		if len(params) < 1 {
			panic(fmt.Sprintf("First parameter of component must be '$' (container)"))
		}

		// 2. find container page/component
		container := params[0].(core.Protoner)

		// 3. create lifecircle controler
		lcc := lifecircle.NewComponentFlow(container, result.Segment.Proton, params[1:])
		lcc.Flow()
		handleComponentReturn(lcc, result.Segment)
		return template.HTML(lcc.String)
	}
}

// handle component return
func handleComponentReturn(lcc *lifecircle.LifeCircleControl, seg *register.ProtonSegment) {
	// no error, no templates return or redirect.
	if seg != nil && lcc.Err == nil && lcc.ResultType == "" {
		// find default tempalte to return
		identity, templatePath := seg.TemplatePath()
		// debug.Log("-756- [ComponentTemplateSelect] %v -> %v", key, tplPath)
		if _, err := templates.GotTemplateCache.Get(identity, templatePath); err != nil {
			lcc.Err = err
		} else {
			// fmt.Println("render component tempalte " + key)
			var buffer bytes.Buffer
			if err := templates.RenderGotTemplate(&buffer, identity, lcc.Proton); err != nil {
				lcc.Err = err
			}
			lcc.String = buffer.String()
		}
	}
	if lcc.Err != nil {
		panic(lcc.Err.Error())
	}
}
