/*
   Time-stamp: <[lifecircle-page.go] Elivoa @ Saturday, 2013-07-27 00:28:03>
*/
package lifecircle

import (
	"fmt"
	"got/core"
	"got/utils"
	"net/http"
	"reflect"
	"strings"
)

// --------------------------------------------------------------------------------

func NewPageFlow(w http.ResponseWriter, r *http.Request, page core.Pager) *LifeCircleControl {
	// init & maintaince structCache
	if si := scache.GetPageX(reflect.TypeOf(page)); si == nil {
		panic("Can't parse page!")
	}
	return newLifeCircleControl(w, r, core.PAGE, page)
}

var pageLifecircles = []string{
	"Setup",
	"SetupRender", // SetupRender is deprecated. use Setup instead.
	"BeginRender",
	"AfterRender",
}

// Page Render Flow:       new -> path -> url ->
func (lcc *LifeCircleControl) Flow() *LifeCircleControl {

	// Inject
	lcc.injectBasic().injectPath().injectURLParameter()
	// lcc.InjectValue() // old version

	// Acitvate() in Tapestry5 is used to receive parameters in path.
	// I use Tag `path:"#"` to do the same thing. So Activate() receives no parameters.
	// Note: Only Page has Activate() event.
	//       Activate() will be called before an event call on page's event or any inner component's event.
	//
	if lcc.Kind == core.PAGE {
		if ret := lcc.CallEvent("Activate"); ret {
			return lcc
		}
	}

	// On form submit, go to another flow.
	// Note: Now only support post to page.
	//       TODO: support to submit to component.
	//
	if lcc.R.Method == "POST" && lcc.Kind == core.PAGE {
		return lcc.PostFlow()
	}

	// Call PageFlow Events
	for _, eventName := range pageLifecircles {
		if ret := lcc.CallEvent(eventName); ret {
			return lcc
		}
	}

	// TODO Handle Returns Here
	return lcc
}

// --------  Event Call on Page  -------------------------------------------------------------------

func (lcc *LifeCircleControl) EventCall(event string) *LifeCircleControl {

	// Note that page is new created. all values needs inject.

	// 1. Inject values into root page
	lcc.injectBasic().injectPath().injectURLParameter()

	// 2. Call Activate method on root page
	if lcc.Kind == core.PAGE {
		if ret := lcc.CallEvent("Activate"); ret {
			return lcc
		}
	}

	// NOTE: Performance: Todo: It's no need to use proton. A
	// reflect.Type is enough, and thus don't need to create new value
	// to each node on path. But how to get proton.Kind() only use a reflect.Type?

	proton := lcc.Proton                    // proton is near upper container proton
	eventPaths := strings.Split(event, ".") // event call path

	// follow the path
	for _, piece := range eventPaths[0 : len(eventPaths)-1] {
		// 1. get from proton cache;
		// !!!This can't be happened!!!
		// !! This is no use, because every request is a new proton value.
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

		// create new proton instance. maybe a component
		var newProton core.Protoner
		fi := si.FieldInfo(piece)
		if fi != nil { // field info cached.
			newInstance := newInstance(fi.Type)
			newProton = newInstance.Interface().(core.Protoner)
		} else {
			// If not cached fieldInfo, create FieldInfo
			containerType := utils.GetRootType(proton)
			field, ok := containerType.FieldByName(piece) // component in path
			if !ok {
				panic(fmt.Sprintf("Can't get field in path: %v", piece))
			}
			newInstance := newInstance(field.Type)                   // create new instance
			newProton = newInstance.Interface().(core.Protoner)      //
			si.CacheEmbedProton(field.Type, piece, newProton.Kind()) // cache
		}
		//// lcc.InjectValueTo(newProton) // don't inject to passby nodes.
		proton.SetEmbed(piece, newProton) // store newInstance into proton
		proton = newProton                // next round
	}

	// inject into new component. make it's value available.
	lcc.injectBasicTo(proton)
	lcc.injectPathTo(proton)
	lcc.injectURLParameterTo(proton)

	// last node, call the event.
	piece := eventPaths[len(eventPaths)-1]
	// Call event. TODO add parameters.
	fmt.Println("\n----------    EVENT CALL    ----------------")
	fmt.Printf("Call event [%v] with parameters.\n", event)
	if ret := lcc._callEventWithURLParameters("On"+piece, reflect.ValueOf(proton)); ret {
		return lcc
	}
	return lcc
}
