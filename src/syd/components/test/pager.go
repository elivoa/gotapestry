package test

import (
	"bytes"
	"github.com/elivoa/got/core"
	"strconv"
)

// decorader
var (
	d_pagenumber1 = ""
	d_pagenumber2 = ""
	d_curr1       = "<"
	d_curr2       = ">"
	d_spliter     = " | "
	d_prev_label  = "Previous"
	d_next_label  = "Next"
	d_first_label = "First" // not used
	d_last_label  = "Last"  // not used
)

type TestPager struct {
	core.Component

	Total     int
	Current   int
	PageItems int
	HTML      string
}

func (p *TestPager) SetupRender() {
	html, err := p.GeneratePagerHtml()
	if err != nil {
		panic("panic when generating pager html. Error is: " + err.Error())
	}
	p.HTML = html
}

// TODO move to Pager file.
func (p *TestPager) GeneratePagerHtml() (string, error) {
	p.FixData()

	var buffer bytes.Buffer
	left := p.Total
	i := 1
	for left > 0 {
		// generate spliter before
		if i > 1 {
			buffer.WriteString(d_spliter)
		}

		// main page nubmer
		if p.PageItems*(i-1) < p.Current && p.Current <= p.PageItems*i {
			buffer.WriteString(d_curr1)
			buffer.WriteString(strconv.Itoa(i))
			buffer.WriteString(d_curr2)
		} else {
			buffer.WriteString(d_pagenumber1)
			buffer.WriteString(strconv.Itoa(i))
			buffer.WriteString(d_pagenumber2)
		}
		// prepare next loop
		i += 1
		left -= p.PageItems
	}

	return buffer.String(), nil
}

// FixData fix invalid data. return true if has errors and fixed.
func (p *TestPager) FixData() bool {
	var hasError = false
	if p.Total < 0 {
		p.Total = 0
		hasError = true
	}
	if p.Current < 0 {
		p.Current = 0
		hasError = true
	}
	if p.PageItems <= 0 {
		p.PageItems = 10
		hasError = true
	}
	return hasError
}

// Check checks if all the values is right. TODO not used
func (p *TestPager) Check() (bool, error) {
	return false, nil
}
