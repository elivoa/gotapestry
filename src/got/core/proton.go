package core

import (
	"bytes"
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
	IncLoopIndex() int                     // is this useful?
	ClientId() string                      // no meaning for PAGE // TODO: generate clientid in loop
	CID() string                           // short for ClientId()
	SetId(id string)

	// attached *plifecircle.Life
	FlowLife() interface{}        // used to store *lifecircle.Life
	SetFlowLife(life interface{}) // set *lifecircle.Life into page.

	// for componenter
	SetInformalParameters(ips *InformalParameters)
	AddInformalParameter(key string, value interface{}) int
	InformalParameters() *InformalParameters
	InformalParameter(key string) interface{}
	InformalParameterString() string
}

// A Proton is a Page, Component or Mixins.
type Proton struct {
	// buildin
	W http.ResponseWriter
	R *http.Request

	Tid       string // component id
	LoopIndex int    // used when component are in a loop

	injected map[string]bool     // field that successfully injected
	embed    map[string]Protoner // embed components [id -> protoner]
	flowlife interface{}         // should be *lifecircle.Life

	// for component only
	informalParameters *InformalParameters // map[string]interface{}
}

type InformalParameters struct {
	Order []string
	Data  map[string]interface{}
}

func NewInformalParameters() *InformalParameters {
	return &InformalParameters{
		Order: []string{},
		Data:  map[string]interface{}{},
	}
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

func (p *Proton) Embed(name string) (Protoner, bool) {
	proton, ok := p.embed[name]
	fmt.Println("\t&&&&&&&&&&&&&&&&")
	for k, _ := range p.embed {
		fmt.Println("\tProtonn embed:", k)
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

func (p *Proton) IncLoopIndex() int {
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

func (p *Proton) CID() string {
	return p.ClientId()
}

func (p *Proton) SetId(id string) {
	p.Tid = id
}

// SetInformalParameters(ips *InformalParameters)
// AddInformalParameter(key string, value interface{}) int
// InformalParameters() *InformalParameters
// InformalParameter(key string) interface{}

// for componenter
func (c *Proton) SetInformalParameters(ips *InformalParameters) {
	c.informalParameters = ips
}

func (c *Proton) AddInformalParameter(key string, value interface{}) int {
	if nil == c.informalParameters {
		c.informalParameters = NewInformalParameters()
	}
	c.informalParameters.Data[key] = value
	c.informalParameters.Order = append(c.informalParameters.Order, key)
	return len(c.informalParameters.Data)
}

func (c *Proton) InformalParameters() *InformalParameters {
	return c.informalParameters
}

func (c *Proton) InformalParameter(key string) interface{} {
	if nil != c.informalParameters {
		return c.informalParameters.Data[key]
	}
	return nil
}

// TODO ordered.
func (c *Proton) InformalParameterString() string {
	var buffer bytes.Buffer
	if nil != c.informalParameters {
		for key, value := range c.informalParameters.Data {
			buffer.WriteString(fmt.Sprintf("%s=\"%v\" ", key, value))
		}
	}
	return buffer.String()
}

// ----------------------------------------------------------------------------------------------------
// TEST: should be deleted
func (p *Proton) ShowInjected() {
	for k, v := range p.injected {
		fmt.Printf(" %v --> %v\n", k, v)
	}
}
