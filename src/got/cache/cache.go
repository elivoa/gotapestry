/*
  Time-stamp: <[cache.go] Elivoa @ Wednesday, 2014-04-23 14:45:54>
  Cache Page/Component Struct info.
  And Component/mixins neasted info.

  Copyright 2012 The Gorilla Authors. All rights reserved.
  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.

  TODO:
    - print all contents.
*/

package cache

import (
	"encoding/json"
	"fmt"
	"github.com/elivoa/got/parser"
	"got/core"
	"got/debug"
	"got/utils"
	"reflect"
	"strings"
	"sync"
)

// --------------------------------------------------------------------------------
// SourceInfo Cache. from package parser.

var SourceCache *parser.SourceInfo

// --------------------------------------------------------------------------------

var StructCache = NewCache() // PUBLIC

func NewCache() *Cache {
	return &Cache{
		m: make(map[reflect.Type]*StructInfo),
	}
}

// ________________________________________
//
type Cache struct {
	l sync.RWMutex
	m map[reflect.Type]*StructInfo
}

type StructInfo struct {
	l      sync.RWMutex
	t      reflect.Type
	Fields map[string]*FieldInfo // lower(fieldname) -> FieldInfo
	Kind   core.Kind             // current struct's kind
}

// IsSlice()
type FieldInfo struct {
	Name string       // field's name
	Type reflect.Type // field's type, or slice's Elem() type.
	Kind core.Kind    // [page|component|mixin|struct|slice?]

	// common fields
	Index   int  // field's index in the struct, not used yet!
	IsSlice bool // is this fields a slice type? Used in InjectFormValues.

	// embed component
	Tid string // if Kind is component, this is Tid, default component's name
}

// ________________________________________________________________________________
// functions for StructInfo
//

func (si *StructInfo) IsPage() bool {
	return si.Kind == core.PAGE
}

func (si *StructInfo) IsComponent() bool {
	return si.Kind == core.COMPONENT
}

// returns FieldInfo
func (si *StructInfo) FieldInfo(field string) *FieldInfo {
	si.l.RLock()
	fi := si.Fields[strings.ToLower(field)]
	si.l.RUnlock()
	return fi
}

// advanced
func (si *StructInfo) Deep(field string) *StructInfo {
	si.l.RLock()
	fi := si.FieldInfo(field)
	si.l.RUnlock()
	if fi != nil {
		StructCache.l.RLock()
		deeperSI, ok := StructCache.m[fi.Type]
		StructCache.l.RUnlock()
		if ok {
			return deeperSI
		}
	}
	return nil
}

// ________________________________________________________________________________
// functions for Cache;; TODO Redefine this funcs, too bad names.
//

// GetX get cache by type. If not found, add to cache.
func (c *Cache) GetX(rt reflect.Type) *StructInfo {
	return c.GetCreate(rt, core.STRUCT)
}

func (c *Cache) GetPageX(rt reflect.Type) *StructInfo {
	return c.GetCreate(rt, core.PAGE)
}

func (c *Cache) GetComopnentX(rt reflect.Type) *StructInfo {
	return c.GetCreate(rt, core.COMPONENT)
}

// Type must be struct. Other type has no meanings.
func (c *Cache) GetCreate(rt reflect.Type, kind core.Kind) *StructInfo {
	// 1. prepare
	t, _ := utils.RemovePointer(rt, true) // remove ptr and slice of type
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("[got/cache]: %v must be struct.", rt))
	}

	// 2. get and return
	c.l.RLock()
	si := c.m[t]
	c.l.RUnlock()
	if si != nil {
		// validation
		if si.Kind != kind {
			panic(fmt.Sprintf(
				"Cached Struct[%v] with kind[%v] is not desired %v, maybe this is caused by conflict of id.",
				si.t, si.Kind, kind))
		}
		return si
	}

	// 3. if not cached, generate StructInfo
	si = c.create(t, kind)

	// set back
	c.l.Lock()
	c.m[t] = si
	c.l.Unlock()
	return si
}

// creat creates a structInfo with meta-data about a struct.
// For Embed Components: When cache page struct, embed components's kind is unknown.
//     store all field as unknown. Modify Kind when render component(Call CacheEmbedProton)
//
func (c *Cache) create(rt reflect.Type, kind core.Kind) *StructInfo {
	t, _ := utils.RemovePointer(rt, false) // already removed, no need to remove again?
	si := &StructInfo{
		t:      t,
		Fields: make(map[string]*FieldInfo),
		Kind:   kind,
	}
	createFieldInfo(si, t)
	if false {
		fmt.Println("*******  print struct info *******************************************************")
		for k, v := range si.Fields {
			fmt.Println(k, " -> ", v)
		}
	}
	return si
}

// analysis fields in t, set info into si.
func createFieldInfo(si *StructInfo, t reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		alias := fieldAlias(field)
		if alias == "-" { // Ignore this field.
			continue
		}

		//
		if field.Anonymous {
			anonymousType, _ := utils.RemovePointer(field.Type, false)
			createFieldInfo(si, anonymousType)
		}

		ft, isSlice := utils.RemovePointer(field.Type, true)
		if isStruct := ft.Kind() == reflect.Struct; !isStruct {
			// GB: don't check if it's supported.
			/*
				if conv := c.conv[ft]; conv == nil {
					// Type is not supported.
					continue
				}
			*/
		}
		fi := &FieldInfo{
			Name:    field.Name,
			Index:   i,
			Type:    field.Type,
			IsSlice: isSlice,
			Kind:    core.UNKNOWN,
		}
		// TODO here judge if it's a page or component, store the type.

		si.l.Lock()
		si.Fields[strings.ToLower(alias)] = fi
		si.l.Unlock()
	}
}

// Append FieldInfo which describ a component embed in a page or component.
// Use component's as field name. panic if conflict with other field names.
// TODO: support mixins.
// TODO: conflict with id version
// Note: tid must has meaningful value.
//
// Return fieldInfo if cached, modify kind if cached or not cached.
//
func (si *StructInfo) CacheEmbedProton(rt reflect.Type, tid string, kind core.Kind) *FieldInfo {
	debug.Log("-759- [Embed Component] register embed component %v with id %v.", rt, tid)

	t, _ := utils.RemovePointer(rt, false)
	if fi := si.FieldInfo(tid); fi != nil {
		// validate kind
		if fi.Kind == kind {
			return fi
		}
		// panic if type mismatch
		if t != fi.Type {
			panic(fmt.Sprintf("Type mismatch, Conflict of proton's ID %v", tid))
		}
		// pass validation, update field
		if fi.Kind == core.UNKNOWN {
			fi.Kind = kind
		}
		return fi
	}

	// if not cached. create FieldInfo and cache.
	fi := &FieldInfo{
		// Name:    tid, // if component with tid but no field, use tid as name.
		Tid:     tid,
		Kind:    core.COMPONENT,
		Index:   -1, // not exist in proton
		Type:    t,
		IsSlice: false,
	}

	si.l.Lock()
	si.Fields[strings.ToLower(tid)] = fi
	si.l.Unlock()
	return fi
}

// ________________________________________________________________________________
// String them
//

func (c *Cache) String() string {
	c.l.Lock()
	str := []string{
		"\n---  -------- StructCache --------",
		fmt.Sprintf("  +  struct cached %v items.", len(c.m)),
	}
	for k, v := range c.m {
		str = append(str, fmt.Sprintf("   : %-20v --> %v", k, v))
	}
	c.l.Unlock()
	return strings.Join(str, "\n")
}

func (fi *FieldInfo) String() string {
	b, err := json.Marshal(fi)
	if err != nil {
		return fmt.Sprintf("[name:StructInfo{Type:%v, Idx:%v, IsSlice:%v}",
			fi.Type, fi.Index, fi.IsSlice)
	}
	return string(b)
}

// __________________________________________________________________________________________
// helper
//

// fieldAlias parses a field tag to get a field alias.
func fieldAlias(field reflect.StructField) string {
	var alias string
	if tag := field.Tag.Get("schema"); tag != "" {
		// For now tags only support the name but let's folow the
		// comma convention from encoding/json and others.
		if idx := strings.Index(tag, ","); idx == -1 {
			alias = tag
		} else {
			alias = tag[:idx]
		}
	}
	if alias == "" {
		alias = field.Name
	}
	return alias
}
