package register

import (
	"fmt"
	"reflect"
	"sync"
)

// segments

var (
	// lock is not needed here.
	Pages      = ProtonSegment{Name: "/"}
	Components = ProtonSegment{Name: "/"}
	Mixins     = ProtonSegment{Name: "/"}
)

// type maps

var (
	// first time initialized. lock seems not useful for them.
	PageTypeMap      = map[reflect.Type]*ProtonSegment{}
	ComponentTypeMap = map[reflect.Type]*ProtonSegment{}
	MixinTypeMap     = map[reflect.Type]*ProtonSegment{}
)

var TemplateKeyMap = &TempalteKeyMapStruct{
	Keymap: map[string]*ProtonSegment{},
}

type TempalteKeyMapStruct struct {
	l      sync.RWMutex
	Keymap map[string]*ProtonSegment // template key as key
}

// ----------------------------------------------------------------------------------------------------

// Get methods
func GetPage(t reflect.Type) *ProtonSegment {
	if v, ok := PageTypeMap[t]; ok {
		return v
	}
	return nil
}

func GetComponent(t reflect.Type) *ProtonSegment {
	if v, ok := ComponentTypeMap[t]; ok {
		return v
	}
	return nil
}

func GetMixin(t reflect.Type) *ProtonSegment {
	if v, ok := MixinTypeMap[t]; ok {
		return v
	}
	return nil
}

func DeubgPrintTypeMaps() {
	fmt.Println("\n_______________________\nPrint All Registrys by Type:")
	fmt.Println("> pages:")
	for k, v := range PageTypeMap {
		fmt.Printf("\t %s -> %s\n", k, v)
	}
	fmt.Println("> components:")
	for k, v := range ComponentTypeMap {
		fmt.Printf("\t %s -> %s\n", k, v)
	}
	fmt.Println("> mixins:")
	for k, v := range MixinTypeMap {
		fmt.Printf("\t %s -> %s\n", k, v)
	}
}
