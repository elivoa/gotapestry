package core

import (
	"fmt"
	"net/http"
)

// GOT Kind

type Kind uint

const (
	INVALID Kind = iota
	PAGE
	COMPONENT
	MIXIN
	STRUCT
)

/*_______________________________________________________________________________
  Proton
*/
type Protoner interface {
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	Kind() Kind // [page|component|mixin]
}

// Common object which Page and Component both has.
type Proton struct {
	// buildin
	W http.ResponseWriter
	R *http.Request

	injected   map[string]bool      // field that successfully injected
	components map[string]*Protoner // embed components TODO
}

func (p *Proton) Request() *http.Request {
	return p.R
}

func (p *Proton) ResponseWriter() http.ResponseWriter {
	return p.W
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

func (p *Proton) SetDefault(fieldName string, value interface{}) {

}

// TEST: should be deleted
func (p *Proton) ShowInjected() {
	for k, v := range p.injected {
		fmt.Printf(" %v --> %v\n", k, v)
	}
}
