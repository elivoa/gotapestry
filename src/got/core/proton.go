package core

import (
	"fmt"
	"net/http"
	"reflect"
)

// GOT Kind

type Kind uint

const (
	UNKNOWN Kind = iota // invalid
	PAGE
	COMPONENT
	MIXIN
	STRUCT // Normal Struct
)

var KindLabels = map[Kind]string{
	UNKNOWN:   "Unknown",
	PAGE:      "Page",
	COMPONENT: "Component",
	MIXIN:     "Mixin",
	STRUCT:    "struct",
}

// Protoner is interface of Proton
type Protoner interface {
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	SetRequest(*http.Request)
	SetResponseWriter(http.ResponseWriter)
	Kind() Kind
	Injected(fieldName string) bool
	SetInjected(fieldName string, b bool)
	Embed(name string) (Protoner, bool)
	SetEmbed(name string, proton Protoner) // return loop index
	IncEmbed() int
	ClientId() string // no meaning for PAGE
	SetId(id string)

	// attached lifecircle.Life
	FlowLife() interface{}
	SetFlowLife(life interface{})
}

// A Proton is a Page, Component or Mixins.
type Proton struct {
	// buildin
	W http.ResponseWriter
	R *http.Request

	Tid       string // component id
	LoopIndex int    // used when component are in a loop

	injected map[string]bool     // field that successfully injected
	embed    map[string]Protoner // embed components TODO

	flowlife interface{} // should be *lifecircle.Life
}

func (p *Proton) FlowLife() interface{} {
	return p.flowlife
}

func (p *Proton) SetFlowLife(life interface{}) {
	p.flowlife = life
}

func (p *Proton) Request() *http.Request {
	return p.R
}

func (p *Proton) ResponseWriter() http.ResponseWriter {
	return p.W
}

func (p *Proton) SetRequest(r *http.Request) {
	p.R = r
	p.SetInjected("R", true)
}

func (p *Proton) SetResponseWriter(w http.ResponseWriter) {
	p.W = w
	p.SetInjected("W", true)
}

func (p *Proton) Kind() Kind {
	return UNKNOWN
}

// if value is injected by got.
// e.g.: to distingush 0 and NaN when parse param.
func (p *Proton) Injected(fieldName string) bool {
	_, ok := p.injected[fieldName]
	return ok
}

// called by got framework
func (p *Proton) SetInjected(fieldName string, b bool) {
	if p.injected == nil {
		p.injected = make(map[string]bool)
	}
	p.injected[fieldName] = b
}

// func (p *Proton) SetDefault(fieldName string, value interface{}) {

// }

// TEST: should be deleted
func (p *Proton) ShowInjected() {
	for k, v := range p.injected {
		fmt.Printf(" %v --> %v\n", k, v)
	}
}

func (p *Proton) Embed(name string) (Protoner, bool) {
	proton, ok := p.embed[name]
	fmt.Println("&&&&&&&&&&&&&&&&")
	for k, _ := range p.embed {
		fmt.Println("embed:", k)
	}
	return proton, ok
}

func (p *Proton) SetEmbed(name string, proton Protoner) {
	if p.embed == nil {
		p.embed = make(map[string]Protoner)
	}
	_, ok := p.embed[name]
	if ok {
		fmt.Println(reflect.TypeOf(proton))
		panic(fmt.Sprintf("Conflict Embed Component '%v'", name))
	}
	p.embed[name] = proton
	proton.SetId(name)
	proton.SetInjected("Tid", true)
}

func (p *Proton) IncEmbed() int {
	p.LoopIndex += 1
	return p.LoopIndex
}

func (p *Proton) ClientId() string {
	if !p.Injected("Tid") {
		panic("Call ClientId() before Tid be injected!")
	}
	if p.LoopIndex == 0 {
		return p.Tid
	} else {
		return fmt.Sprintf("%v_%v", p.Tid, p.LoopIndex)
	}
}

func (p *Proton) SetId(id string) {
	p.Tid = id
}
