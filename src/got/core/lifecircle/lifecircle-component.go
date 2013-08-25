/*
   Time-stamp: <[lifecircle-component.go] Elivoa @ Saturday, 2013-08-24 14:49:24>
*/
package lifecircle

import (
	"fmt"
	"github.com/gorilla/context"
	"got/core"
	"got/debug"
	"got/register"
	"got/templates"
	"got/utils"
	"html/template"
	"log"
	"path"
	"reflect"
	"strings"
)

// ComponentLifeCircle returns template-func to handle component render.
func ComponentLifeCircle(name string) func(...interface{}) interface{} {

	// returnd string or template.HTML are inserted into final template.
	return func(params ...interface{}) interface{} {

		log.Printf("-620- [flow] Render Component %v ....", name)

		// 1. find base component type
		result, err := register.Components.Lookup(name)
		{
			if err != nil || result.Segment == nil {
				panic(fmt.Sprintf("Component %v not found!", name))
			}
			if len(params) < 1 {
				panic(fmt.Sprintf("First parameter of component must be '$' (container)"))
			}
		}

		// 2. find container page/component
		container := params[0].(core.Protoner)

		// get lcc from component
		{
			fmt.Println("---- debug =============================")
			// fmt.Println(container)
			fmt.Println(container.Request())
			// fmt.Println(context.Get(container.Request(), LCC_OBJECT_KEY))
		}
		lcc := context.Get(container.Request(), LCC_OBJECT_KEY).(*LifeCircleControl)
		life := lcc.componentFlow(container, result.Segment.Proton, params[1:])
		life.SetRegistry(result.Segment)

		// templates renders in common flow()
		returns := life.flow()
		if returns.breakReturn() {
			lcc.returns = returns
			lcc.rendering = false
			// here don't process returns, handle return in page-flow's end.
			// here only set returns into control and stop the rendering.
		}

		// If returns is not template-renderer (i.e.: redirect or text output),
		// flow breaks and will not reach here.
		// Here returns default template render.
		return template.HTML(life.out.String())
	}
}

// --------------------------------------------------------------------------------
//
// Create a new Component Flow.
// param:
//   container - real container object.
//   component - current component base object.
//   params - parameters in the component grammar.
//
// Note: I maintain StructCache here in the flow create func. This occured only when
//       page or component are rendered. Directly post to a page can not invoke structcache init.
//
// TODO: Performance Improve to Component in Loops.
//
func (lcc *LifeCircleControl) componentFlow(container core.Protoner, componentSeed core.Componenter, params []interface{}) *Life {

	{
		debuglog("----- [Create Component flowcontroller] ------------------------%v",
			"----------------------------------------")
		debug.Log("- C - [Component Container] Type: %v, ComponentType:%v,\n",
			reflect.TypeOf(container), reflect.TypeOf(componentSeed))
	}

	// Store type in StructCache, Store instance in ProtonObject.
	// Warrning: What if use component in page/component but it's not initialized?
	// Tid= xxx in template must the same with fieldname in .go file.
	//

	// 1. cache in StructInfoCache. (application scope)
	si := scache.GetCreate(reflect.TypeOf(container), container.Kind())
	if si == nil {
		panic(fmt.Sprintf("StructInfo for %v can't be null!", reflect.TypeOf(container)))
	}
	t := utils.GetRootType(componentSeed)
	tid, _ := determinComponentTid(params, t)
	si.CacheEmbedProton(t, tid, componentSeed.Kind())

	// 2. store in proton's embed field. (request scope)
	proton, ok := container.Embed(tid)
	if !ok {
		// first: create and append.
		life := lcc.appendComponent(componentSeed)
		// proton = life.Proton.(core.Componenter)
		container.SetEmbed(tid, life.proton)
	} else {
		// already found. maybe this component is in a loop or range.
		lcc.current.out.Reset() // components in loop is one instance.
		proton.IncEmbed()
	}
	lcc.injectBasic()
	lcc.injectComponentParameters(params) // inject component parameters
	return lcc.current
}

// return (name, is setManually); t must not be ptr.
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

// --------------------------------------------------------------------------------

// flow controls the common lifecircles, including pages and components.
func (l *Life) flow() (returns *Returns) {
	// There are 2 way to reach here.
	// 1. Page lifecircle, from PageFlow()
	// 2. Component's template-func, from func call. Get lcc from Request.

	// Here follows the flow of tapestry:
	//   http://tapestry.apache.org/component-rendering.html
	//
	// TODO: call lifecircle events with parameter

	for {
		returns = eventReturn(l.call("Setup", "SetupRender"))
		if returns.breakReturn() {
			return
		}
		if !returns.returnsFalse() {

			for {
				returns = eventReturn(l.call("BeginRender"))
				if returns.breakReturn() {
					return
				}
				if !returns.returnsFalse() {

					for {
						returns = eventReturn(l.call("BeforeRenderTemplate"))
						if returns.breakReturn() {
							return
						}
						if !returns.returnsFalse() {

							// Here we ignored BeforeRenderBody and AfterRenderBody.
							// Maybe add it later.
							// May be useful for Loop component?
							// TODO here render template:
							l.renderTemplate()

							// if any component breaks it's render, stop all rendering.
							if l.control.rendering == false {
								returns = nil
								return
							}
						}

						returns = eventReturn(l.call("AfterRenderTemplate"))
						if returns.breakReturn() {
							return
						}
						if !returns.returnsFalse() {
							break
						}
					}
				}
				returns = eventReturn(l.call("AfterRender"))
				if returns.breakReturn() {
					return
				}
				if !returns.returnsFalse() {
					break
				}
			}
		}

		returns = eventReturn(l.call("Cleanup", "CleanupRender"))
		if returns.breakReturn() {
			return
		}
		if !returns.returnsFalse() {
			break // exit
		}
	}

	// finally I go through all render phrase.
	returns = &Returns{
		returnType: "template",
	}
	return
}

// renderTemplate find and render Template using go way.
func (l *Life) renderTemplate() {
	// reach here means I can find the template and render it.
	// I can panic if template not found.
	// debug.Log("-755- [TemplateSelect] %v -> %v", identity, templatePath)

	identity, templatePath := l.registry.TemplatePath()
	if _, err := templates.Cache.Get(identity, templatePath); err != nil {
		panic(err.Error())
	}
	if err := templates.RenderTemplate(&l.out, identity, l.proton); err != nil {
		panic(err.Error()) // lcc.Err = err
	}
}
