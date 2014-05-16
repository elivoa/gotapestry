/*
   Time-stamp: <[lifecircle-component.go] Elivoa @ Saturday, 2014-05-17 00:59:51>
*/
package lifecircle

import (
	"fmt"
	"github.com/elivoa/got/logs"
	"github.com/elivoa/got/route/exit"
	"github.com/elivoa/got/templates"
	"got/core"
	"got/register"
	"got/utils"
	"html/template"
	"path"
	"reflect"
	"strings"
)

var cflog = logs.Get("ComponentFLow")

// ComponentLifeCircle returns template-func to handle component render.
// Param: name e.g.: layout/gotheader
func ComponentLifeCircle(name string) func(...interface{}) interface{} {

	// returnd string or template.HTML are inserted into final template.
	// params: 1. container, 2. uniqueID, 3. other parameter pairs...
	return func(params ...interface{}) interface{} {

		if cflog.Info() {
			cflog.Printf("[Component] Render Component %v (ID:%s) ....", name, params[1].(string))
		}

		// Processing component in one page
		// 1. find component by component's type
		// 2. find component's container
		// 3. create component's life
		// 4. process returns.

		// 1. find base component type
		result, err := register.Components.Lookup(name)
		if err != nil || result.Segment == nil {
			panic(fmt.Sprintf("Component %v not found!", name))
		}
		if len(params) < 1 {
			panic(fmt.Sprintf("First parameter of component must be '$' (container)"))
		}

		// 2. find container page or component
		container := params[0].(core.Protoner)
		tid := params[1].(string)
		containerLife := container.FlowLife().(*Life)
		{
			// fmt.Println("~~~~~~==~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			// fmt.Println("container:", utils.GetRootType(container), "  >>> ", container.Kind())
			// fmt.Println("seed comp:", reflect.TypeOf(result.Segment.Proton))
			// // fmt.Println("component:", lcc.current.rootType, " >>> ", result.Segment.Proton.Kind())
			// // fmt.Println("container:", utils.GetRootType(container), "  >>> ", container.ClientId())
			// // fmt.Println("seed comp:", reflect.TypeOf(result.Segment.Proton))
			// // fmt.Println("component:", lcc.current.rootType, " >>> ", result.Segment.Proton.ClientId())
			// fmt.Println("\n")
		}
		// unused: get lcc from component; use method to get from controler.
		// lcc := context.Get(container.Request(), config.LCC_OBJECT_KEY).(*LifeCircleControl)
		lcc := containerLife.control // container's lcc has the same R and W.
		life := lcc.componentFlow(containerLife, result.Segment.Proton, tid, params[2:])
		life.SetRegistry(result.Segment)

		// templates renders in common flow()
		returns := life.flow()
		if returns.IsBreakExit() {
			lcc.returns = returns
			lcc.rendering = false
			// here don't process returns, handle return in page-flow's end.
			// here only set returns into control and stop the rendering.
		}

		// If returns is not template-renderer (i.e.: redirect or text output),
		// flow breaks and will not reach here.
		// Here returns default template render.
		rr := template.HTML(life.out.String())

		return rr
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
func (lcc *LifeCircleControl) componentFlow(containerLife *Life, componentSeed core.Componenter, tid string, params []interface{}) *Life {

	fmt.Println("\n--------------------------------------------------------------------------------")
	fmt.Println("tid is: ", tid)
	for _, v := range params {
		fmt.Println("param: ", v)
	}

	if cflog.Debug() {
		cflog.Printf("----- [Component flow] ------------------------%v",
			"----------------------------------------")
		cflog.Printf("- C - [Component Container] Type: %v, ComponentType:%v, tid:%s\n",
			containerLife.rootType, reflect.TypeOf(componentSeed), tid)
	}

	// Store type in StructCache, Store instance in ProtonObject.
	// Warrning: What if use component in page/component but it's not initialized?
	// Tid= xxx in template must the same with fieldname in .go file.
	//

	// 1. cache in StructInfoCache. (application scope)
	si := scache.GetCreate(containerLife.rootType, containerLife.kind)
	if si == nil {
		panic(fmt.Sprintf("StructInfo for %v can't be null!", containerLife.rootType))
	}

	// TODO: Is below useful? Can i remove these codes?

	t := utils.GetRootType(componentSeed)
	// tid, _ := determinComponentTid(params, t)
	si.CacheEmbedProton(t, tid, componentSeed.Kind()) // TODO: is this useful?

	// 2. store in proton's embed field. (request scope)
	// containerLife := container.FlowLife().(*Life)
	life, found := containerLife.embedmap[tid]
	if !found {
		// create component life!
		life := containerLife.appendComponent(componentSeed, tid)
		containerLife.proton.SetEmbed(tid, life.proton)
	} else {
		// already exist, must be in a loop.
		lcc.current = life
		// already found. maybe this component is in a loop or range.
		lcc.current.out.Reset()    // components in loop is one instance.
		life.proton.IncLoopIndex() // increase loop index
	}

	lcc.injectBasicTo(lcc.current.proton)
	lcc.injectComponentParameters(params) // inject component parameters to current life
	return lcc.current
}

// return (name, is setManually); t must not be ptr.
func determinComponentTid(params []interface{}, t reflect.Type) (tid string, specifiedId bool) {
	for idx, p := range params {
		if idx%2 == 0 {
			key := strings.ToLower(p.(string))
			if key == "tid" || key == "t:id" {
				tid = params[idx+1].(string)
				// use specified id, can't be duplicated (judge in transform.go)
				specifiedId = true
				return
			}
		}
	}
	// if fall here, measn no specified id.
	// TODO: get id
	tid = path.Ext(t.String())[1:]
	return
}

// --------------------------------------------------------------------------------

// flow controls the common lifecircles, including pages and components.
func (l *Life) flow() (returns *exit.Exit) {
	// There are 2 way to reach here.
	// 1. Page lifecircle, from PageFlow()
	// 2. Component's template-func, from func call. Get lcc from Request.

	// Here follows the flow of tapestry:
	//   http://tapestry.apache.org/component-rendering.html
	//
	// TODO: call lifecircle events with parameter
	// TODO: flat them
	for {
		returns = SmartReturn(l.call("Setup", "SetupRender"))
		if returns.IsBreakExit() {
			return
		}
		if !returns.IsReturnsFalse() {

			for {
				returns = SmartReturn(l.call("BeginRender"))
				if returns.IsBreakExit() {
					return
				}
				if !returns.IsReturnsFalse() {

					for {
						returns = SmartReturn(l.call("BeforeRenderTemplate"))
						if returns.IsBreakExit() {
							return
						}
						if !returns.IsReturnsFalse() {

							// Here we ignored BeforeRenderBody and AfterRenderBody.
							// Maybe add it later.
							// May be useful for Loop component?
							l.renderTemplate()

							// if any component breaks it's render, stop all rendering.
							if l.control.rendering == false {
								returns = nil
								return
							}
						}

						returns = SmartReturn(l.call("AfterRenderTemplate"))
						if returns.IsBreakExit() {
							return
						}
						if !returns.IsReturnsFalse() {
							break
						}
					}
				}
				returns = SmartReturn(l.call("AfterRender"))
				if returns.IsBreakExit() {
					return
				}
				if !returns.IsReturnsFalse() {
					break
				}
			}
		}

		returns = SmartReturn(l.call("Cleanup", "CleanupRender"))
		if returns.IsBreakExit() {
			return
		}
		if !returns.IsReturnsFalse() {
			break // exit
		}
	}

	// finally I go through all render phrase.
	returns = exit.Template(nil)
	return
}

// renderTemplate find and render Template using go way.
func (l *Life) renderTemplate() {
	// reach here means I can find the template and render it.
	// debug.Log("-755- [TemplateSelect] %v -> %v", identity, templatePath)
	if _, err := templates.LoadTemplates(l.registry, false); err != nil {
		panic(err)
	}
	if err := templates.RenderTemplate(&l.out, l.registry.Identity(), l.proton); err != nil {
		panic(err)
	}
}
