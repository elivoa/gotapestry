package route

import (
	"fmt"
	"github.com/gorilla/mux"
	"got/core/lifecircle"
	"got/debug"
	"got/register"
	"got/templates"
	"log"
	"net/http"
	"reflect"
	rd "runtime/debug"
)

var (
	emptyParameters = []reflect.Value{}
	debugLog        = true
)

// ________________________________________________________________________________
// GOT Tapestry style Handler
//
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
// -------- Simple Handler --------
// --------------------------------------------------------------------------------

// func RedirectHandler(url string) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		http.Redirect(w, r, url, http.StatusFound)
// 	}
// }

// // register each page and cache them.
// func PageHandler(basePage core.Pager) func(http.ResponseWriter, *http.Request) {
// 	log.Printf("[building] Init page '%v'", reflect.TypeOf(basePage))

// 	return func(w http.ResponseWriter, r *http.Request) {
// 		printAccessHeader(r)
// 		lcc := lifecircle.NewPageFlow(w, r, basePage).Flow()
// 		handleReturn(lcc, nil)
// 		printAccessFooter(r)
// 	}
// }
