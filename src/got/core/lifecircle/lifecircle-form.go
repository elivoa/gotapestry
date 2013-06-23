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
*/
func (lcc *LifeCircleControl) InjectFormValues() {

	if debug.FLAG_print_form_submit_details && lcc.Kind == "page" {
		debug.PrintFormMap("~ 1 ~ Request.Form", lcc.R.Form)
	}

	// 为了迎合gorilla/schema的奇葩要求，这里需要转换FormData
	// version 1: for form.keys for path.segments.
	// TODO version 2: l2cache

	data := map[string][]string{}

	v := lcc.V
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		lcc.Err = errors.New("got/lifecircle: interface must be a pointer to struct")
		return
	}
	// -----------------------------------------------------------------------
	// parse array in form
	for path, formValue := range lcc.R.Form {
		fmt.Printf("  ~~ processing key:path:%v v.type: %v\n", path, v.Type())

		// is need translate
		fmt.Printf("++++ lcc:%v\n", lcc)
		template, needTranslate := isNeedTranslate(path, v.Type())
		if needTranslate {
			data[path] = formValue
		} else {
			for idx, value := range formValue {
				key := fmt.Sprintf(template, idx)
				data[key] = []string{value}
			}
		}
	}

	if debug.FLAG_print_form_submit_details && lcc.Kind == "page" {
		debug.PrintFormMap("~ 2 ~ gorilla/schema Data", data)
	}

	// decode
	utils.SchemaDecoder.Decode(lcc.Proton, data)

	if debug.FLAG_print_form_submit_details {
		fmt.Printf("++++++++++ lcc.Proton = %v=n", lcc.Proton)
		fmt.Println("\n+END FORM SUBMIT LOG+ <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n")
	}
}

// --------------------------------------------------------------------------------

var tcache = NewTranslateCache()

// top cache of value path
type translateCache struct {
	l sync.Mutex
	m map[reflect.Type]*translateInfo
}

// m: map[path]template - template='' means noneed to translate
type translateInfo struct {
	l sync.Mutex
	m map[string]string
}

func NewTranslateCache() *translateCache {
	return &translateCache{m: make(map[reflect.Type]*translateInfo)}
}

func (c *translateCache) GetOrInitInfo(t reflect.Type) *translateInfo {
	c.l.Lock()
	ti := c.m[t]
	if ti == nil {
		ti = &translateInfo{m: make(map[string]string)}
		c.m[t] = ti
	}
	c.l.Unlock()
	return ti
}

// create info, cache, and return tempate.
// TODO: for now, only support 1 slice objects in it.
func (i *translateInfo) Create(path string, t reflect.Type) string {
	var typo = t

	pieces := strings.Split(path, ".")
	segs := make([]string, 0)
	hasSlice := false
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Printf("++ CREATE TRANSLATE_INFO: PATH:%v on Type:%v --> %v\n", path, t, pieces)
	for _, p := range pieces {
		// adjust typo to the write struct type.
		allow, isSlice := false, false

		fmt.Println("++++++++  ------------------------------------------  +++++++++++++++++")
		fmt.Printf("++++ DDDDDDDDDDDDDD: path:%-6v piece:%-10v  TYPO:%v\n", path, p, typo)
		fmt.Printf("++++ DDDDDDDDDDDDDD: typo.Kind() = %v\n", typo.Kind())
		// if typo.Kind() == reflect.Ptr {
		// 	fmt.Println("typo is ptr")
		// }

		switch typo.Kind() {
		case reflect.Ptr:
			elem := typo.Elem()
			switch elem.Kind() {
			case reflect.Struct:
				allow = true
			case reflect.Slice:
				isSlice = true
			default:
				allow = false
			}
			typo = typo.Elem() // remove pointer
		case reflect.Struct:
			allow = true
		case reflect.Slice:
			isSlice = true
		default:
			allow = false
		}

		if isSlice {
			allow = true
			typo = typo.Elem()        // remove slice pointer
			segs = append(segs, "%d") // append number place-holder

			// set global flag
			if hasSlice { // forbbid 2 slice in one path.
				panic("lifecircle: Don't allow 2 level slice elements. " + path)
			}
			hasSlice = true
		}

		// get StructInfo if is struct value
		// set next round type to fields' type. nil if not go on.
		var structInfo *cache.StructInfo = nil
		if allow {
			structInfo = c.GetnCache(typo)
			fieldInfo, ok := structInfo.Fields[p]
			if !ok { // no such field, stop.
				typo = nil
			} else {
				typo = fieldInfo.Type
			}
		} else {
			typo = nil
		}

		// append the 1
		segs = append(segs, p)
	}

	var template string = ""
	if hasSlice {
		template = strings.Join(segs, ".")
	}

	i.l.Lock()
	i.m[path] = template
	i.l.Unlock()

	//	fmt.Printf(" *********** create template for %v is: %v\n", path, template)
	return template
}

// v.Kind() must be struct
// return template, isNeedTranslate?
func isNeedTranslate(path string, t reflect.Type) (string, bool) {
	ti := tcache.GetOrInitInfo(t)

	var (
		template string
		ok       bool
	)
	ti.l.Lock()
	template, ok = ti.m[path]
	ti.l.Unlock()
	if !ok {
		template = ti.Create(path, t)
	}
	return template, template == ""
}

func translateRequestFromIntoGorillaForm(
	requestForm map[string][]string) *map[string][]string {
	return nil
}
