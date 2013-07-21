/*
  Time-stamp: <[cache.go] Elivoa @ Sunday, 2013-07-21 12:10:16>
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
	"fmt"
	"got/core"
	"got/debug"
	"got/utils"
	"reflect"
	"strings"
	"sync"
)

var StructCache = NewCache() // PUBLIC

func NewCache() *Cache {
	return &Cache{
		m: make(map[reflect.Type]*StructInfo),
	}
}

// ________________________________________
//
type Cache struct {
	l sync.Mutex
	m map[reflect.Type]*StructInfo
}

type StructInfo struct {
	l      sync.Mutex
	t      reflect.Type
	Fields map[string]*FieldInfo
	Kind   core.Kind // [page|component|mixin|struct] current struct's kind
}

// IsSlice()
type FieldInfo struct {
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

func (si *StructInfo) FieldInfo(field string) *FieldInfo {
	si.l.Lock()
	fi := si.Fields[field]
	si.l.Unlock()
	return fi
}

// advanced
func (si *StructInfo) Deep(field string) *StructInfo {
	if fi := si.FieldInfo(field); fi != nil {
		// TODO mutex?
		StructCache.l.Lock()
		deeperSI, ok := StructCache.m[fi.Type]
		StructCache.l.Unlock()
		if ok {
			return deeperSI
		}
	}
	return nil
}

// ________________________________________________________________________________
// functions for Cache
//
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
	c.l.Lock()
	si := c.m[t]
	c.l.Unlock()
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
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		alias := fieldAlias(field)
		if alias == "-" { // Ignore this field.
			continue
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
			Index:   i,
			Type:    field.Type,
			IsSlice: isSlice,
			Kind:    core.UNKNOWN, // unknown if it's that.
		}
		// TODO here judge if it's a page or component, store the type.

		si.l.Lock()
		si.Fields[alias] = fi
		si.l.Unlock()
	}
	return si
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
		Tid:     tid,
		Kind:    core.COMPONENT,
		Index:   -1, // not exist in proton
		Type:    t,
		IsSlice: false,
	}

	si.l.Lock()
	si.Fields[tid] = fi
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
	return fmt.Sprintf("StructInfo{Type:%v, Idx:%v, IsSlice:%v}",
		fi.Type, fi.Index, fi.IsSlice)
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
