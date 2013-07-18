// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*

*/

package cache

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var StructCache = NewCache()

func NewCache() *Cache {
	return &Cache{
		m: make(map[reflect.Type]*StructInfo),
	}
}

type Cache struct {
	l sync.Mutex
	m map[reflect.Type]*StructInfo
}

type StructInfo struct {
	t      reflect.Type
	Fields map[string]*FieldInfo
}

type FieldInfo struct {
	Type    reflect.Type
	Idx     int  // field idnex in the struct, what is this used for?
	IsSlice bool // is this fields a slice type? used in InjectFormValues.
}

// ________________________________________________________________________________
// functions
//

func (fi *FieldInfo) String() string {
	return fmt.Sprintf("StructInfo{Type:%v, Idx:%v, IsSlice:%v}",
		fi.Type, fi.Idx, fi.IsSlice)
}

// type must be struct. others has no meanings.
func (c *Cache) GetnCache(rt reflect.Type) *StructInfo {
	t := rt
	// remove ptr
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// remove slice
	if t.Kind() == reflect.Slice {
		t = t.Elem()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		// fmt.Printf("finally i found slice type: %v\n", t)
	}
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("[got/cache]: %v must be struct.", rt))
	}

	// get and return
	c.l.Lock()
	si := c.m[t]
	c.l.Unlock()
	if si != nil {
		return si
	}

	// not cached, generate StructInfo
	si = c.create(t)
	// fmt.Printf(".... [CACHE] Create StructInfo for [%v]\n", t)

	// set back
	c.l.Lock()
	c.m[t] = si
	c.l.Unlock()
	return si
}

// creat creates a structInfo with meta-data about a struct.
func (c *Cache) create(t reflect.Type) *StructInfo {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	info := &StructInfo{
		t:      t,
		Fields: make(map[string]*FieldInfo),
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		alias := fieldAlias(field)
		if alias == "-" {
			// Ignore this field.
			continue
		}

		if false { // DEBUG PRINT
			fmt.Printf(".... [cache] field: %v (%v)......\n",
				field, alias,
			)
		}

		// Check if the type is supported and don't cache it if not.
		// First let's get the basic type.
		isSlice, isStruct := false, false
		ft := field.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		// fmt.Printf(">>>>> %v\n", ft)
		if isSlice = ft.Kind() == reflect.Slice; isSlice {
			ft = ft.Elem()
			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}
		}
		// fmt.Printf(">>>>> %v\n", ft)
		if isStruct = ft.Kind() == reflect.Struct; !isStruct {
			// GB: don't check if it's supported.
			/*
				if conv := c.conv[ft]; conv == nil {
					// Type is not supported.
					continue
				}
			*/
		}
		info.Fields[alias] = &FieldInfo{
			Idx:     i,
			Type:    field.Type,
			IsSlice: isSlice, // && isStruct,
		}
	}
	return info
}

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

// --------------------------------------------------------------------------------

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
