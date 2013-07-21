package register

import (
	"bytes"
	"fmt"
	"got/core"
	"got/core/lifecircle"
	"got/debug"
	"got/templates"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

/* ________________________________________________________________________________
   ComponentRegister
*/
var Components = ProtonSegment{Name: "/"}

func Component(f func(), components ...core.Componenter) int {
	for _, c := range components {
		url := makeUrl(f, c)
		selectors := Components.Add(url, c, "component")

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
  Handler method
  Return: string or template.HTML
*/
func componentLifeCircle(name string) func(...interface{}) interface{} {

	return func(params ...interface{}) interface{} {

		log.Printf("-620- [flow] Render Component %v ....", name)

		// 1. find base component type
		result, err := Components.Lookup(name)
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
func handleComponentReturn(lcc *lifecircle.LifeCircleControl, seg *ProtonSegment) {
	// no error, no templates return or redirect.
	if seg != nil && lcc.Err == nil && lcc.ResultType == "" {
		// find default tempalte to return
		key, tplPath := LocateGOTComponentTemplate(seg.Src, seg.Path)
		// debug.Log("-756- [ComponentTemplateSelect] %v -> %v", key, tplPath)
		_, err := templates.GotTemplateCache.Get(key, tplPath)
		if nil != err {
			lcc.Err = err
		} else {
			// fmt.Println("render component tempalte " + key)
			var buffer bytes.Buffer
			err = templates.RenderGotTemplate(&buffer, key, lcc.Proton)
			if err != nil {
				lcc.Err = err
			}
			lcc.String = buffer.String()
		}
	}

	if lcc.Err != nil {
		debug.Error(lcc.Err)
		http.Error(lcc.W, fmt.Sprint(lcc.Err), http.StatusInternalServerError)
	}
}

// ________________________________________________________________________________
// Locate Templates
// return (template-key, template-file-path); TODO: performance issue
func LocateGOTComponentTemplate(src string, path string) (string, string) {
	appConfig := Apps.Get(src)
	if appConfig == nil {
		panic(fmt.Sprintf("Can't find APP Config %v", src))
	}
	key := fmt.Sprintf("c_%v:%v", src, path)
	templateFilePath := filepath.Join(appConfig.FilePath, "components", path) + ".html"
	return key, templateFilePath
}
