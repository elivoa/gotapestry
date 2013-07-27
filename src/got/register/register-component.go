package register

import (
	"bytes"
	"fmt"
	"got/core"
	"got/core/lifecircle"
	"got/templates"
	"html/template"
	"log"
	"strings"
)

/* ________________________________________________________________________________
   ComponentRegister
*/
var Components = ProtonSegment{Name: "/"}

func Component(f func(), components ...core.Componenter) int {
	for _, c := range components {
		url := makeUrl(f, c)
		selectors := Components.Add(url, c)
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
