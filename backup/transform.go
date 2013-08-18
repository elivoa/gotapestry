/**
  Time-stamp: <[transform.go] Elivoa @ Friday, 2013-08-16 17:51:19>
*/
package transform

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"errors"
	"fmt"
	"got/cache"
	"got/core"
	"got/register"
	"io"
	"reflect"
	"regexp"
	"strings"
)

// ---- Transform template ------------------------------------------
type Transformater struct {
	tree *Node
	b    bytes.Buffer
	z    *html.Tokenizer
}

func NewTransformer() *Transformater {
	return &Transformater{}
}

/*
  Transform tempalte fiels. functions:
  translate <t:some_component ... /> into {{t_some_component ...}}

TODO:
  . Support t:block
  . Range Tag

TODOs:
---- 1 --------------------------------------------------------------------------------
<div t:type="xx"... >some <b>bold text</b></div>
 there will remaining: some meaningful text</div>
 now I ignore these, TODO make this a block and render it.

---- N --------------------------------------------------------------------------------
*/
var compressHtml bool = false

func (t *Transformater) Parse(reader io.Reader) *Transformater {
	z := html.NewTokenizer(reader)
	t.z = z
	for {
		tt := z.Next()
		// after call something all tag is lowercased. but here with case.
		zraw := z.Raw()
		raw := make([]byte, len(zraw))
		copy(raw[:], zraw[:])
		var ()
		switch tt {
		case html.TextToken:
			// here may contains {{ }}
			// trim spaces?
			if compressHtml {
				t.b.Write(TrimTextNode(z.Raw())) // trimed spaces
			} else {
				t.b.Write(raw)
			}
		case html.StartTagToken:
			if b := t.processStartTag(); !b {
				t.b.Write(raw)
			}
		case html.SelfClosingTagToken:
			if b := t.processStartTag(); !b {
				t.b.Write(raw)
			}
		case html.EndTagToken:
			k, _ := z.TagName()
			switch string(k) {
			case "range", "with", "if":
				t.b.WriteString("{{end}}")
			case "hide":
				t.b.WriteString("*/}}")
			default:
				t.b.Write(raw)
			}
		// case html.CommentToken:
		// 	// ignore all comments
		// // case html.DoctypeToken:
		case html.ErrorToken:
			if z.Err().Error() == "EOF" {
				return t
			} else {
				panic(z.Err().Error())
			}
		default:
			t.b.Write(raw)
		}
	}
	return t
}

// processing every start tag()
// return
//   - true if already write to buffer.
//   - false if need to write Raw() to buffer.
// Note: go.net/html package lowercased all values,
//
func (t *Transformater) processStartTag() bool {
	// collect information
	bname, hasAttr := t.z.TagName()
	var (
		iscomopnent   bool
		componentName []byte
		elementName   []byte
		err           error
	)
	if len(bname) >= 2 && bname[0] == 't' && bname[1] == ':' {
		iscomopnent = true
		componentName = bname[2:]
	}

	var attrs map[string][]byte
	// var attrs [][][]byte
	if hasAttr {
		attrs = map[string][]byte{}
		for {
			key, val, more := t.z.TagAttr()
			if len(key) == 6 && bytes.Equal(key, []byte("t:type")) {
				iscomopnent = true
				componentName = val
				elementName = bname
				// ignore t:type attr
			} else {
				attrs[string(key)] = val // = append(attrs, [][]byte{key, val})
			}
			if !more {
				break
			}
		}
	}
	if iscomopnent {
		if err = t.transformComponent(componentName, elementName, attrs); err == nil {
			return true
		}
	}

	// --------------------------------------------------------------------------------
	// not a component, process if tag is command
	switch string(bname) {
	case "range":
		t.renderRange(attrs)
	case "if":
		t.renderIf(attrs)
	case "else":
		t.b.WriteString("{{else}}")
	case "hide":
		t.b.WriteString("{{/*")
	case "t:import":
		t.b.WriteString("----------")
	default:
		if err != nil {
			panic(err.Error())
		}
		return false
	}
	return true
}

func (t *Transformater) renderRange(attrs map[string][]byte) {
	t.b.WriteString("{{range ")
	if nil != attrs {
		if _var, ok := attrs["var"]; ok {
			t.b.Write(_var)
			t.b.WriteString(":=")
		}
		if source, ok := attrs["source"]; ok {
			t.b.Write(source)
		}
	}
	t.b.WriteString("}}")
}

func (t *Transformater) renderIf(attrs map[string][]byte) {
	t.b.WriteString("{{if ")
	if nil != attrs {
		var (
			_var []byte
			ok   bool
		)
		if _var, ok = attrs["t"]; !ok {
			if _var, ok = attrs["test"]; !ok {
				panic("`If` must have attribute test or t!")
			}
		}
		t.b.Write(_var)
	}
	t.b.WriteString("}}")
}

func (t *Transformater) transformComponent(componentName []byte, elementName []byte,
	attrs map[string][]byte) error {

	// lookup component and get StructInfo
	lookupurl := strings.Replace(string(componentName), ".", "/", -1)
	lr, err := register.Components.Lookup(lookupurl)
	if err == nil && (lr.Segment == nil || lr.Segment.Proton == nil) {
		err = errors.New(fmt.Sprintf("Can't find component for %v", string(componentName)))
	}
	if err != nil {
		return err
	}

	sc := cache.StructCache
	si := sc.GetCreate(reflect.TypeOf(lr.Segment.Proton), core.COMPONENT)
	// TODO: cache embed directly elements.

	t.b.WriteString("{{t_")
	// t.b.Write(componentName)
	t.b.WriteString(strings.Replace(lookupurl, "/", "_", -1))
	t.b.WriteString(" $")

	// elementName
	if elementName != nil {
		t.b.WriteString(" \"elementName\" `")
		t.b.Write(elementName)
		t.b.WriteString("`")
	}

	if attrs != nil {
		for key, val := range attrs {
			// write key, all capitlize
			t.b.WriteString(" \"")
			// get which is cached.
			fi := si.FieldInfo(key)
			if fi != nil {
				t.b.WriteString(fi.Name)
			} else {
				t.b.WriteString(key)
			}
			t.b.WriteString("\" ")

			// TODO: Auto-detect literal or functional
			// Value transform: for name="_some_value_", we transform it into:
			//   ~ before ~           ~ after ~             ~ note ~
			//   ".Name"              .Name                // start form . or $
			//   "literal:....."      "...."               // literal prefix
			//	 "abcd"               "abcd"               // auto detect plan text
			//	 ".Name+'_'+.Id"      (print .Name '_' .Id)// special
			//   "/xxx/{{.ID}}"       (print "/xxx/" .Id)
			//
			//   TODO support more prefix...
			//
			// if value starts from . or $ , treate this as property. others as string
			switch {
			case len(val) > 0 && (val[0] == '.' || val[0] == '$' || val[0] == '('):
				t.b.Write(val)
			case len(val) > 5 && bytes.Equal(val[0:5], []byte("print")):
				t.b.WriteString("(")
				t.b.Write(val)
				t.b.WriteString(")")
			case len(val) > 8 && bytes.Equal(val[0:8], []byte("literal:")):
				t.b.WriteString(" \"")
				t.b.Write(bytes.Replace(val[8:], []byte{'"'}, []byte{'\\', '"'}, 0))
				t.b.WriteString("\"")
			case printValueRegex.Match(val): // if is "/xxx/{{.ID}}"
				result := printValueRegex.FindSubmatch(val)
				// for _, r := range result {
				// 	fmt.Println(r)
				// }
				if len(result) == 3 { // translate to (print "/xxx/" .ID)
					t.b.WriteString(" (print \"")
					t.b.Write(result[1])
					t.b.WriteString("\" ")
					t.b.Write(result[2])
					t.b.WriteString(")")
				}
			default:
				t.b.WriteString(" \"")
				t.b.Write(bytes.Replace(val, []byte{'"'}, []byte{'\\', '"'}, 0))
				t.b.WriteString("\"")
			}
		}
	}
	t.b.WriteString("}}")
	return nil
}

var printValueRegex, _ = regexp.Compile("^(.*){{(.*)}}$")

func (t *Transformater) Render() string {
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println(t.b.String())
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++=+++")
	return t.b.String()
}

// ---- utils --------------------------------------------------------------------------------

// TODO trim node function not finished.
func TrimTextNode(text []byte) []byte {
	// var (
	// 	addSpaceLeft    bool = false
	// 	addSpaceRight   bool = false
	// 	addNewLineLeft  bool = false
	// 	addNewLineRight bool = false
	// )
	// for _, b := range bytes {
	// 	// if b
	// }
	return bytes.Trim(text, " ")
}
