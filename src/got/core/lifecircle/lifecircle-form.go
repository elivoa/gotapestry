/*
Lifecircle-form

Inject values when submit.

*/

package lifecircle

import (
	"errors"
	"fmt"
	"got/cache"
	"got/debug"
	"got/utils"
	"reflect"
	"strings"
	"sync"
)

// quick access
var c = cache.StructCache

/*
  Parse values submited from Form. (use gorilla/schema)
  TODO:
    . Support parse repeated values.
    . Support File upload

  BUG:
*/
func (lcc *LifeCircleControl) InjectFormValues() {

	// add ParseForm to fix bugs in go1.1.1
	err := lcc.R.ParseForm()
	if err != nil {
		lcc.Err = err
		return
		// panic(err.Error())
	}

	// debug print
	if debug.FLAG_print_form_submit_details && lcc.Kind == "page" {
		debug.PrintFormMap("~ 1 ~ Request.Form", lcc.R.Form)
	}

	// 为了迎合gorilla/schema的奇葩要求，这里需要转换格式为：FormData
	// version 1: for form.keys for path.segments.
	// TODO version 2: l2cache

	// 1) Precondition
	v := lcc.V
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		lcc.Err = errors.New("got/lifecircle: interface must be a pointer to struct")
		return
	}

	// ---------------------------------
	// 2) Parse array in form
	data := map[string][]string{} // stores transfered FormData

	// don't import multipart form.
	// err = lcc.R.ParseMultipartForm(1024 * 1024 * 10) // MOVE CONFIG OUT
	// if err != nil {
	// 	panic(err.Error())
	// }

	for path, formValue := range lcc.R.Form {

		// ------------------------------------------------------------
		// Get something from cache.
		// is path in name Attribute need translate to "x.y.1.z" format
		var (
			leafType reflect.Type
			template string
			ok       bool
		)
		ti := tcache.getTranslateInfo(v.Type())
		ti.l.Lock()
		template, ok = ti.templates[path]
		ti.l.Unlock()
		if !ok {
			template, leafType = ti.Create(path, v.Type())
		} else {
			leafType = ti.types[path]
		}
		// ---- END ----

		{ // DEBUG PRINT
			var k string
			if leafType == nil {
				k = "nil"
			} else {
				k = leafType.Kind().String()
			}
			fmt.Printf(" ~~ processing key-path [%-20v],"+
				" leafKind:[%v], template:[%-20v] v.type: %v\n",
				path, k, template, v.Type(),
			)
		} // ---- END DEBUG ----

		// issue #4 in github.com/gorilla/schema
		// this is just a fix. filter out all empty string.
		if leafType != nil && leafType.Kind() == reflect.Slice {
			switch leafType.Elem().Kind() {
			case reflect.Int, reflect.Float32, reflect.Float64, reflect.Int64:
				formValue = _killEmpty(formValue, "0")
			case reflect.String:
				formValue = _killEmpty(formValue, " ")
			}
		}

		if template == "" {
			// no need to translate, copy the value
			data[path] = formValue
		} else {
			for idx, value := range formValue {
				key := fmt.Sprintf(template, idx)
				data[key] = []string{value}
			}
		}
	}

	// debug print
	if debug.FLAG_print_form_submit_details && lcc.Kind == "page" {
		debug.PrintFormMap("~ 2 ~ gorilla/schema Data", data)
	}

	// 3) decode
	utils.SchemaDecoder.Decode(lcc.Proton, data)

	// debug print
	if debug.FLAG_print_form_submit_details {
		fmt.Printf("++++++++++ lcc.Proton = %v=n", lcc.Proton)
		fmt.Println("\n+END FORM SUBMIT LOG+ <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<" +
			"<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n")
	}
}

func _killEmpty(slices []string, emptyTo string) []string {
	nvalue := make([]string, len(slices))
	// empty value to 0
	for idx, value := range slices {
		if strings.Trim(value, " ") == "" {
			nvalue[idx] = emptyTo
		} else {
			nvalue[idx] = value
		}
	}
	return nvalue
}

func translateRequestFromIntoGorillaForm(
	requestForm map[string][]string) *map[string][]string {
	return nil
}

// ________________________________________________________________________________
// The Cache part. L2Cache (L1Cache is got/cache)
//
var tcache = NewTranslateCache()

func NewTranslateCache() *translateCache {
	return &translateCache{m: make(map[reflect.Type]*translateInfo)}
}

//
// Top Translate Cache stores which type should translate.
//
type translateCache struct {
	l sync.Mutex
	m map[reflect.Type]*translateInfo
}

// translateInfo stores the [translate-path] of a [path].
// template-path = '' means no need to translate
type translateInfo struct {
	l         sync.Mutex
	templates map[string]string       // path -> template
	types     map[string]reflect.Type // path -> last node type
}

// get translateInfo from cache, if nil, create an empty one.
func (c *translateCache) getTranslateInfo(t reflect.Type) *translateInfo {
	c.l.Lock()
	ti := c.m[t]
	if ti == nil {
		// init new translateInfo
		ti = &translateInfo{
			templates: make(map[string]string),
			types:     make(map[string]reflect.Type),
		}
		c.m[t] = ti
	}
	c.l.Unlock()
	return ti
}

/*
  Create template and last-node-type for `path` in translateInfo.
  cache them then return.
  Param:
    path - select path;
    t - root type
  Return
  TODO:
    for now, only support 1 slice objects in it.
*/
func (i *translateInfo) Create(path string, t reflect.Type) (string, reflect.Type) {
	var parentType = t
	pieces := strings.Split(path, ".")
	segs := make([]string, 0)
	hasSlice := false

	// debug print
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Printf("++ CREATE TRANSLATE_INFO: PATH:%v on Type:%v --> %v\n", path, t, pieces)

	// last node's type, if slice, it's slice's element type.
	// NOTE:!! for now only return slice element's type.
	var leafType reflect.Type

	// 1. start loop in path
	for idx, p := range pieces {

		// debug print
		// {
		// 	fmt.Println("++++++++  ------------------------------------------  ++++++++")
		// 	fmt.Printf("++++ DDDDDDDDDDDDDD: path:%-6v piece:%-10v  parentType:%v\n",
		// 		path, p, parentType)
		// 	fmt.Printf("++++ DDDDDDDDDDDDDD: typo.Kind() = %v\n", parentType.Kind())
		// 	fmt.Println("....................................................")
		// }

		// 3. append path segments to template.
		//   ** root can't be slice type
		segs = append(segs, p)

		// 1. get StructInfo from cache.
		structInfo := c.GetnCache(parentType)
		if nil == structInfo {
			panic("struct info is null for " + parentType.String())
		}

		fieldInfo, ok := structInfo.Fields[p]
		if ok && fieldInfo != nil {
			leafType = fieldInfo.Type
			if fieldInfo.IsSlice {
				// if leafe node && is slice, stop here.
				if idx == len(pieces)-1 {
					break
				}
				hasSlice = true
				segs = append(segs, "%d") // append number place-holder
			} else if fieldInfo.Type.Kind() == reflect.Struct {
				// continue
			} else {
				// stop here
			}
		}

		//
		if ok && fieldInfo != nil {
			// fmt.Printf("---- set parentType to fieldInfo.Type: %v\n", fieldInfo.Type)
			// fmt.Printf("---- fieldinfo is: %v\n", fieldInfo)
			parentType = fieldInfo.Type
		} else {
			// no such field, stop.
			parentType = nil
			// fmt.Printf("++++ FieldInfo for path-segment[%v] is nil, !!so break here!!\n", p)
			break
		}
	}

	// construct template
	var template string = ""
	if hasSlice {
		template = strings.Join(segs, ".")
	}

	// set to TranslateInfo
	i.l.Lock()
	i.templates[path] = template
	i.types[path] = leafType
	i.l.Unlock()

	fmt.Printf("**** Create template for %v is: %v\n", path, template)
	fmt.Printf("**** Final type is: %v\n\n", leafType)
	return template, leafType

}
