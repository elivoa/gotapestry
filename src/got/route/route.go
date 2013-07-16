package route

import (
	"fmt"
	"github.com/gorilla/mux"
	"got/core"
	"got/core/lifecircle"
	"got/debug"
	"got/register"
	"got/templates"
	"log"
	"net/http"
	"path/filepath"
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

	// if error occured
	printAccessHeader(r)
	defer func() {
		if err := recover(); err != nil {
			processPanic(err)
			// TODO render 500 page.
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		}
		printAccessFooter(r)
	}()

	// 1. let's find the right pages.
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

	debug.Log("-601- [RouteFind] %v", result.Segment)

	// go through lifecircle
	lcc := lifecircle.NewPageFlow(w, r, result.Segment.Proton)
	lcc.SetPageUrl(result.PageUrl).SetEventName(result.EventName)
	if result.EventName == "" {
		lcc.Flow() // normal flow
		handleReturn(lcc, result.Segment)
	} else {
		lcc.EventCall(result.EventName) // event call
		// when call event, process common return ('redirect', "...")
		// TODO 1: default return the current page.
	}
}

func handleReturn(lcc *lifecircle.LifeCircleControl, seg *register.ProtonSegment) {
	// no error, no templates return or redirect.
	if seg != nil && lcc.Err == nil && lcc.ResultType == "" {
		// find default tempalte to return
		key, tplPath := LocateGOTTemplate(seg.Src, seg.Path)
		// debug.Log("-755- [TemplateSelect] %v -> %v", key, tplPath)
		_, err := templates.GotTemplateCache.Get(key, tplPath)
		if nil != err {
			lcc.Err = err
		} else {
			// fmt.Println("render tempalte " + key)
			e := templates.RenderGotTemplate(lcc.W, key, lcc.Proton)
			if e != nil {
				lcc.Err = e
			}
		}
	}

	if lcc.Err != nil {
		panic(lcc.Err.Error())
	}
}

// ________________________________________________________________________________
// Locate Templates

// performance issue
// return (template-key, template-file-path)
func LocateGOTTemplate(src string, path string) (string, string) {
	// println("\n >> locate template")
	appConfig := register.Apps.Get(src)
	if appConfig == nil {
		panic(fmt.Sprintf("Can't find APP Config %v", src))
	}
	key := fmt.Sprintf("%v:%v", src, path)
	templateFilePath := filepath.Join(appConfig.FilePath, "pages", path) + ".html"
	return key, templateFilePath
}

// --------------------------------------------------------------------------------
// -------- Simple Handler --------
// --------------------------------------------------------------------------------

func RedirectHandler(url string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, http.StatusFound)
	}
}

// register each page and cache them.
func PageHandler(basePage core.IPage) func(http.ResponseWriter, *http.Request) {
	log.Printf("[building] Init page '%v'", reflect.TypeOf(basePage))

	return func(w http.ResponseWriter, r *http.Request) {
		printAccessHeader(r)
		lcc := lifecircle.NewPageFlow(w, r, basePage).Flow()
		handleReturn(lcc, nil)
		printAccessFooter(r)
	}
}

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
	fmt.Println("-----------------------------^         PAGE RENDER END           -----------------------------------")
	fmt.Println("....................................................................................................")
}

func processPanic(err interface{}) {
	log.Printf("xxxxxxxx  PANIC  xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	log.Printf("x panic: %-109v x", err)
	log.Printf("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	fmt.Println("StackTrace >>")
	rd.PrintStack()
	fmt.Println()
}
