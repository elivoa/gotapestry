// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cache

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var StructCache = NewCache()

func NewCache() *Cache {
	c := Cache{
		m: make(map[reflect.Type]*StructInfo),
	}
	return &c
}

type Cache struct {
	l sync.Mutex
	m map[reflect.Type]*StructInfo
}

// type must be struct. others has no meanings.
func (c *Cache) GetnCache(rt reflect.Type) *StructInfo {
	t := rt
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("got/cache: must be struct")
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
	fmt.Printf("....... [cache] Create StructInfo: %v \n", t)

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
		// Check if the type is supported and don't cache it if not.
		// First let's get the basic type.
		isSlice, isStruct := false, false
		ft := field.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		if isSlice = ft.Kind() == reflect.Slice; isSlice {
			ft = ft.Elem()
			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}
		}
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
			IsSlice: isSlice && isStruct,
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

// StructInfo
type StructInfo struct {
	t      reflect.Type
	Fields map[string]*FieldInfo
}

// FieldInfo
type FieldInfo struct {
	Type    reflect.Type
	Idx     int  // field idnex in the struct, what is this used for?
	IsSlice bool // is this fields a slice type? used in InjectFormValues.
}

func (fi *FieldInfo) String() string {
	return fmt.Sprintf("StructInfo{Type:%v, Idx:%v, IsSlice:%v}",
		fi.Type, fi.Idx, fi.IsSlice)
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
