/*
   Time-stamp: <[lifecircle-page.go] Elivoa @ Sunday, 2014-04-20 12:51:31>
*/
package lifecircle

import (
	"fmt"
	"github.com/elivoa/got/config"
	"got/core"
	"got/register"
	"got/utils"
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
	lcc.injectBasic().injectPath().injectURLParameter()

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

func (lcc *LifeCircleControl) EventCall(event string) *LifeCircleControl {

	// Note that page is new created. all values needs inject.

	// 1. Inject values into root page
	lcc.injectBasic().injectPath().injectURLParameter()

	returns := SmartReturn(lcc.page.call("Activate"))
	if !returns.IsReturnsTrue() {
		return lcc
	}

	// NOTE: Performance: Todo: It's no need to use proton. A
	// reflect.Type is enough, and thus don't need to create new value
	// to each node on path. But how to get proton.Kind() only use a reflect.Type?

	proton := lcc.page.proton               // proton is upper container's proton
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
