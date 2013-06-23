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
)

const (
	TemplateKey = "template-key"
	PageKey     = "page-key"
)

var (
	emptyParameters = []reflect.Value{}
	debugLog        = true
)

// ________________________________________________________________________________
// GOT Tapestry style Handler

func RouteHandler(w http.ResponseWriter, r *http.Request) {
	url := "/" + mux.Vars(r)["url"]

	// 1. skip special resources. TODO Expand to config.
	if url == "/favicon.ico" {
		return
	}

	printAccessHeader(r)

	// 1. let's find the write pages.
	seg, pageUrl, err := register.Pages.Lookup(url)
	fmt.Printf(">>>>>>ｓｅｇ　ｉｓ　: %v\n", seg)
	if nil != err {
		panic(err.Error())
	}
	if seg == nil {
		panic(fmt.Sprintf("Error: seg.Proton is null. seg: %v", seg))
	}
	if seg.Proton == nil {
		// TODO redirect to 404 page.
		fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
		panic(fmt.Sprintf("~~~~ Page not found ~~~~"))
	}

	debug.Log("-601- [RouteFind] %v", seg)

	// go through lifecircle
	lcc := lifecircle.NewPageFlow(w, r, seg.Proton).SetPageUrl(pageUrl).Flow()
	handleReturn(lcc, seg)

	printAccessFooter(r)
}

func handleReturn(lcc *lifecircle.LifeCircleControl, seg *register.ProtonSegment) {
	// no error, no templates return or redirect.
	if seg != nil && lcc.Err == nil && lcc.ResultType == "" {
		// find default tempalte to return
		key, tplPath := LocateGOTTemplate(seg.Src, seg.Path)
		debug.Log("-755- [TemplateSelect] %v -> %v", key, tplPath)
		_, err := templates.GotTemplateCache.Get(key, tplPath)
		if nil != err {
			lcc.Err = err
		} else {
			fmt.Println("render tempalte " + key)
			e := templates.RenderGotTemplate(lcc.W, key, lcc.Proton)
			if e != nil {
				lcc.Err = e
			}
		}
		// var err error
		// if lcc.Kind == "component" {
		// 	var buffer bytes.Buffer
		// 	err = templates.RenderTemplate(&buffer, templateName, lcc.Proton)
		// 	lcc.String = buffer.String()
		// } else {
		// 	err = templates.RenderTemplate(lcc.W, templateName, lcc.Proton)
		// }
	}

	if lcc.Err != nil {
		debug.Error(lcc.Err)
		http.Error(lcc.W, fmt.Sprint(lcc.Err), http.StatusInternalServerError)
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

	// if templateName == "" {
	// 	lcc.Err = errors.New(fmt.Sprintf(
	// 		"%v %v must has an associated template or other return.",
	// 		lcc.Kind, lcc.Name,
	// 	))
	// 	return lcc.Err
	// }

	// debuglog("-980- Render Tempalte %v.", templateName)

	//	return "", ""
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
	log.Printf("\n____________________________________________________________________________________")
	log.Printf(">>> access %v\n", r.URL.Path)
	// log.Printf("> w is %v\n", reflect.TypeOf(w))
	// log.Printf("> w is %v\n", reflect.TypeOf(req))
}

func printAccessFooter(r *http.Request) {
	//debug.Log("^ ^ ^ ^ ^ ^ ^ ^ PAGE RENDER END ^ ^ ^ ^ ^ ^ ^ ^ ^ ^")
	debug.Log("^.^.^.^.^.^.^.^ PAGE RENDER END ^.^.^.^.^ nmn^o^ ^.^.^")
}
