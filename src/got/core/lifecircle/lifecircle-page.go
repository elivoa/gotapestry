/*
   Time-stamp: <[lifecircle-page.go] Elivoa @ Saturday, 2014-05-17 20:52:57>
*/
package lifecircle

import (
	"fmt"
	"github.com/elivoa/got/config"
	"github.com/elivoa/got/logs"
	"github.com/elivoa/got/templates"
	"got/register"
	"net/http"
	"reflect"
	"strings"
)

// --------------------------------------------------------------------------------

func NewPageFlow(w http.ResponseWriter, r *http.Request, registry *register.ProtonSegment) *LifeCircleControl {
	// init & maintaince structCache
	page := registry.Proton
	if si := scache.GetPageX(reflect.TypeOf(page)); si == nil {
		panic("Can't parse page!")
	}
	lcc := newControl(w, r)
	lcc.createPage(page)
	lcc.page.SetRegistry(registry)
	return lcc
}

// Page Render Flow: (Entrance 1)    new -> path -> url ->
func (lcc *LifeCircleControl) PageFlow() *LifeCircleControl {

	// Inject
	lcc.injectBasic().injectPath().injectURLParameter().injectHiddenThings()

	// add lcc object to request.
	lcc.SetToRequest(config.LCC_OBJECT_KEY, lcc)

	// Acitvate() in Tapestry5 is used to receive parameters in path.
	// I use Tag `path:"#"` to do the same thing. Here Activate() receives no parameters.
	// Note: Only Page has Activate() event.
	//       Activate() will also be called before an event call.
	//
	returns := SmartReturn(lcc.page.call("Activate"))
	if returns.IsReturnsTrue() {
		if lcc.r.Method == "POST" {
			// >> Form post flow
			// Note: Now only support post to page.
			//       TODO: support to submit to component.
			returns = lcc.PostFlow()
			if returns.IsBreakExit() {
				lcc.HandleBreakReturn()
			}
		} else {
			// >> page render flow

			// Save lcc in request scope. There are two approaches:
			//   1. First is set lcc to request data store.(now)
			//   2. The other way is set to the proton object. component can get $ object.
			lcc.SetToRequest(config.LCC_OBJECT_KEY, lcc)

			// universial flow
			lcc.rendering = true
			returns = lcc.page.flow()
			lcc.rendering = false
		}
	}

	// error handling
	if lcc.Err != nil {
		panic(lcc.Err.Error())
	}

	if returns == nil {
		// if embed components has break-returns.
		returns = lcc.returns
	} else {
		// set returns back to lcc.
		lcc.returns = returns
	}

	// handle returns
	if returns.IsBreakExit() {
		lcc.HandleBreakReturn() // handle break return.
	} else if !returns.IsReturnsFalse() {
		// normal template-rendering
		lcc.w.Write(lcc.page.out.Bytes())
	}
	return lcc
}

// ----  Event Call on Page  --------------------------------------------------
var eventlog = logs.Get("GOT:EventCall")

// Note: for EventCall result.segment is page's semgnet
// 暂不支持Event的popup
func (lcc *LifeCircleControl) EventCall(result *register.LookupResult) *LifeCircleControl {

	fmt.Println("--$$------------------------------------ Call event ----------------------------")
	// Note that page is new created. all values needs inject.

	// 1. Inject values into root page
	if eventlog.Debug() {
		eventlog.Printf("[EventCall] Trigger EventCall on %v", lcc.r.URL.Path)
	}
	lcc.injectBasic().injectPath().injectURLParameter().injectHiddenThings()

	// call rootpage's activate method before call event.
	returns := SmartReturn(lcc.page.call("Activate"))
	if !returns.IsReturnsTrue() {
		return lcc
	}

	// NOTE: Performance: Todo: It's no need to use proton. A
	// reflect.Type is enough, and thus don't need to create new value
	// to each node on path. But how to get proton.Kind() only use a reflect.Type?

	if result.ComponentPaths != nil && len(result.ComponentPaths) > 0 {
		// Call event on embed components, need create new instance.
		if eventlog.Debug() {
			eventlog.Printf("[EventCall] follow by components: %v", result.ComponentPaths)
		}

		currentSeg := FollowComponentByIds(lcc.page.registry, result.ComponentPaths)
		lcc.current = newLife(currentSeg.Proton) // new instance

		// inject again, because current is changed.
		if eventlog.Debug() {
			eventlog.Printf("[EventCall] Inject things into new leafe node.")
		}
		lcc.injectBasic().injectPath().injectURLParameter().injectHiddenThings()
	} else {
		// Call event on page.
	}

	fmt.Println("\n----------    EVENT CALL    ----------------")
	if eventlog.Debug() {
		eventlog.Printf("[EventCall] Call Event %v", "On"+result.EventName)
	}

	if ret := lcc._callEventWithURLParameters("On"+result.EventName, result.Parameters, lcc.current.v); ret {
		return lcc
	}
	return lcc

	// find components in it. segInHell is the leaf node of component who contains the event.
	// seed := currentSeg.Proton

	// // follow the path
	// for _, piece := range eventPaths[0 : len(eventPaths)-1] {
	// 	// 1. get from proton cache;
	// 	// !!!This can't be happened!!!
	// 	// !! This is no use, because every request is a new proton value.
	// 	c, ok := proton.Embed(piece)
	// 	if ok && c != nil {
	// 		proton = c
	// 		continue
	// 	}

	// 	// 2. Is Cached in StructInfo
	// 	si := scache.GetCreate(reflect.TypeOf(proton), proton.Kind()) // root page
	// 	if si == nil {
	// 		panic(fmt.Sprintf("StructInfo for %v can't be null!", reflect.TypeOf(proton)))
	// 	}

	// 	// create new proton instance. maybe a component
	// 	var newProton core.Protoner
	// 	fi := si.FieldInfo(piece)
	// 	if fi != nil { // field info cached.
	// 		newInstance := newInstance(fi.Type)
	// 		newProton = newInstance.Interface().(core.Protoner)
	// 	} else {
	// 		// If not cached fieldInfo, create FieldInfo
	// 		containerType := utils.GetRootType(proton)
	// 		field, ok := containerType.FieldByName(piece) // component in path
	// 		if !ok {
	// 			panic(fmt.Sprintf("Can't get field in path: %v", piece))
	// 		}
	// 		newInstance := newInstance(field.Type)                   // create new instance
	// 		newProton = newInstance.Interface().(core.Protoner)      //
	// 		si.CacheEmbedProton(field.Type, piece, newProton.Kind()) // cache
	// 	}
	// 	//// lcc.InjectValueTo(newProton) // don't inject to passby nodes.
	// 	proton.SetEmbed(piece, newProton) // store newInstance into proton
	// 	proton = newProton                // next round
	// }

}

func FollowComponentByIds(seg *register.ProtonSegment, componentIds []string) *register.ProtonSegment {
	fmt.Println("\n7788: alsdjflajsdlfj")
	current := seg
	if componentIds != nil {
		for idx, componentId := range componentIds {
			// if component not exists, load and parse it.
			fmt.Printf(">> find %vth component by ID:%v \n", idx, componentId)

			lowercasedId := strings.ToLower(componentId)
			if !current.IsTemplateLoaded {
				fmt.Println("   >> LoadTemplate ", lowercasedId, "")
				if _, err := templates.LoadTemplates(current, false); err != nil {
					panic(err)
				}
			}

			fmt.Println("\n----------------------------------------------------")
			fmt.Println("idx: ", idx, "; componentId", componentId)
			fmt.Println("leng of embed is: ", len(current.EmbedComponents))
			if s, ok := current.EmbedComponents[lowercasedId]; ok {
				// fmt.Println("   >> go into: ", lowercasedId, " >> seg is: ", s)
				current = s
			} else {
				panic(fmt.Sprintf("Can't find component for id:%s ", s))
			}
		}
	}
	return current
}
