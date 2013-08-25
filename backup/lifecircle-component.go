/*
   Time-stamp: <[lifecircle-component.go] Elivoa @ Saturday, 2013-07-27 10:50:48>
*/
package lifecircle

import (
	"fmt"
	"got/core"
	"got/debug"
	"got/utils"
	"path"
	"reflect"
	"strings"
)

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
func NewComponentFlow(container core.Protoner, component core.Componenter,
	params []interface{}) *LifeCircleControl {

	{
		debuglog("----- [Create Component flowcontroller] ------------------------%v",
			"----------------------------------------")
		debug.Log("- C - [Component Container] Type: %v, ComponentType:%v,\n",
			reflect.TypeOf(container), reflect.TypeOf(component))
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
	t := utils.GetRootType(component)
	tid, _ := determinComponentTid(params, t)
	si.CacheEmbedProton(t, tid, component.Kind())

	// 2. store in proton's embed field. (request scope)
	proton, ok := container.Embed(tid)
	var lcc *LifeCircleControl
	if !ok {
		// The first create new component object.
		lcc = newLifeCircleControl(container.ResponseWriter(), container.Request(),
			core.COMPONENT, component)
		container.SetEmbed(tid, lcc.Proton)
	} else {
		// If proton in a loop, we use the same proton instance.
		lcc = &LifeCircleControl{
			W:        container.ResponseWriter(),
			R:        container.Request(),
			Kind:     core.COMPONENT,
			Proton:   proton,
			RootType: utils.GetRootType(proton),
			V:        reflect.ValueOf(proton),
			Name:     fmt.Sprint(reflect.TypeOf(proton).Elem()),
		}
		proton.IncEmbed() // increase loop index. Used by ClientId()
	}

	lcc.injectComponentParameters(params) // inject component parameters
	return lcc
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
