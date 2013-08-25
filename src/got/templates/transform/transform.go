/**
  Time-stamp: <[transform.go] Elivoa @ Saturday, 2013-08-24 14:47:09>
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
	tree   *Node // root node
	blocks map[string]*Node
	z      *html.Tokenizer
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

	// the root node
	root := newNode() // &Node{level: 0}
	t.tree = root
	parent := root

	for {
		tt := z.Next()

		// new the current node.
		node := newNode()

		// after call something all tag is lowercased. but here with case.
		zraw := z.Raw()
		node.raw = make([]byte, len(zraw))
		copy(node.raw, zraw[:])
		zraw = node.raw

		// start parse
		switch tt {
		case html.TextToken:
			// here may contains {{ }}
			if compressHtml {
				node.html.Write(TrimTextNode(z.Raw())) // trimed spaces
			} else {
				node.html.Write(zraw)
			}
			parent.AddChild(node)

		case html.StartTagToken:
			node.closed = false
			if b := t.processStartTag(node); !b {
				node.html.Write(zraw)
			}
			// switch node.tagName {
			// case "input", "br", "hr", "link":
			// 	parent.AddChild(node)
			// default:
			parent.AddChild(node)
			parent = node // go in
			// }

		case html.SelfClosingTagToken:
			if b := t.processStartTag(node); !b {
				node.html.Write(zraw)
			}
			parent.AddChild(node)

		case html.EndTagToken:
			k, _ := z.TagName()
			tag := string(k)
			switch tag {
			case "range", "with", "if":
				node.html.WriteString("{{end}}")
			case "hide":
				node.html.WriteString("*/}}")
			default:
				node.html.Write(zraw)
			}
			// TODO: process unclosed tag.
			// if has unclosed tag, just unclose it.
			// find the right tag and close, move wrong tag back.
			if tag == parent.tagName {
				parent.AddChild(node)
				parent.closed = true
				parent = parent.parent
			} else {
				// fmt.Println(">>>+++++++ ", node)

				node.parent = parent // only set parent will not link the node to the tree.
				temp := node
				for {
					// if true{break}
					if temp == nil {
						panic(fmt.Sprintf("Tag %v not closed!", temp))
					}
					temp = temp.parent
					if tag == temp.tagName {
						temp.AddChild(node)
						parent = temp.parent
						temp.closed = true
						break
					} else {
						if temp.children != nil {
							// fmt.Println("    > ++++++++++++++++++ move children up!", temp)
							// tp := []*Node{}
							for _, c := range temp.children {
								// fmt.Println("      > move <<< ", c.tagName, ";", c.html.String(), ">>>")
								c.Detach()
								temp.parent.AddChild(c)
								temp.closed = true
							}
						}
					}
				}
			}

		// case html.CommentToken:
		// 	// ignore all comments
		// case html.DoctypeToken:
		// 	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++", z.Raw())
		case html.ErrorToken:
			if z.Err().Error() == "EOF" {

				// here is the entrance
				// get blocks from tree.
				t.parseBlocks()

				return t
			} else {
				panic(z.Err().Error())
			}
		// case html.DoctypeToken:
		// 	node.html.Write(zraw)
		default:
			node.html.Write(zraw)
			parent.AddChild(node)
		}
	}
	// can't be here
	return t
}

func (t *Transformater) parseBlocks() {
	t.blocks = map[string]*Node{}
	t._parseBlocks(t.tree)
}

func (t *Transformater) _parseBlocks(n *Node) {
	if n == nil {
		return
	}
	// TODO do something
	if n.tagName == "t:block" {
		var found bool = false
		var id string
		if n.attrs != nil {
			for k, v := range n.attrs {
				if strings.ToLower(k) == "id" {
					id = string(v)
					found = true
					break
				}
			}
		}
		if !found {
			panic("Can't find `id` attribute in t:block tag!")
		}

		t.blocks[id] = n.Detach()
	} else {
		if n.children != nil {
			for _, node := range n.children {
				t._parseBlocks(node)
			}
		}
	}

}

// processing every start tag()
// return 1.
//   - true if already write to buffer.
//   - false if need to write Raw() to buffer.
//   2. tagNamep
// Note: go.net/html package lowercased all values,
//
//
func (t *Transformater) processStartTag(node *Node) bool {
	// collect information
	bname, hasAttr := t.z.TagName()
	node.tagName = string(bname) // performance
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
		node.attrs = attrs
	}
	if iscomopnent {
		if err = t.transformComponent(node, componentName, elementName, attrs); err == nil {
			return true
		}
	}

	// --------------------------------------------------------------------------------
	// not a component, process if tag is command
	switch string(bname) {
	case "range":
		t.renderRange(node, attrs)
	case "if":
		t.renderIf(node, attrs)
	case "else":
		node.html.WriteString("{{else}}")
	case "hide":
		node.html.WriteString("{{/*")
	case "t:import":
		node.html.WriteString("----------")
	case "t:block":
		t.renderBlock(node, attrs)
	case "t:delegate":
		t.renderDelegate(node, attrs)
	default:
		if err != nil {
			panic(err.Error())
		}
		return false
	}
	return true
}

func (t *Transformater) renderBlock(node *Node, attrs map[string][]byte) {
	node.html.WriteString("||delegate some one||")

}

func (t *Transformater) renderDelegate(node *Node, attrs map[string][]byte) {
	node.html.Write(node.raw)
}

func (t *Transformater) builtinComponentFunction(name string) func(*Node, map[string][]byte) {
	switch name {
	case "range":
		return t.renderRange
	default:
		panic(fmt.Sprintf("Builtin component %v not found!", name))
	}
}

func (t *Transformater) renderRange(node *Node, attrs map[string][]byte) {
	node.html.WriteString("{{range ")
	if nil != attrs {
		if _var, ok := attrs["var"]; ok {
			node.html.Write(_var)
			node.html.WriteString(":=")
		}
		if source, ok := attrs["source"]; ok {
			node.html.Write(source)
		}
	}
	node.html.WriteString("}}")
}

func (t *Transformater) renderIf(node *Node, attrs map[string][]byte) {
	node.html.WriteString("{{if ")
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
		node.html.Write(_var)
	}
	node.html.WriteString("}}")
}

func (t *Transformater) transformComponent(node *Node, componentName []byte, elementName []byte,
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

	node.html.WriteString("{{t_")
	// node.html.Write(componentName)
	node.html.WriteString(strings.Replace(lookupurl, "/", "_", -1))
	node.html.WriteString(" $")

	// elementName
	if elementName != nil {
		node.html.WriteString(" \"elementName\" `")
		node.html.Write(elementName)
		node.html.WriteString("`")
	}

	if attrs != nil {
		for key, val := range attrs {
			// write key, all capitlize
			node.html.WriteString(" \"")
			// get which is cached.
			fi := si.FieldInfo(key)
			if fi != nil {
				node.html.WriteString(fi.Name)
			} else {
				node.html.WriteString(key)
			}
			node.html.WriteString("\" ")

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
				node.html.Write(val)
			case len(val) > 5 && bytes.Equal(val[0:5], []byte("print")):
				node.html.WriteString("(")
				node.html.Write(val)
				node.html.WriteString(")")
			case len(val) > 8 && bytes.Equal(val[0:8], []byte("literal:")):
				node.html.WriteString(" \"")
				node.html.Write(bytes.Replace(val[8:], []byte{'"'}, []byte{'\\', '"'}, 0))
				node.html.WriteString("\"")
			case printValueRegex.Match(val): // if is "/xxx/{{.ID}}"
				result := printValueRegex.FindSubmatch(val)
				// for _, r := range result {
				// 	fmt.Println(r)
				// }
				if len(result) == 3 { // translate to (print "/xxx/" .ID)
					node.html.WriteString(" (print \"")
					node.html.Write(result[1])
					node.html.WriteString("\" ")
					node.html.Write(result[2])
					node.html.WriteString(")")
				}
			default:
				node.html.WriteString(" \"")
				node.html.Write(bytes.Replace(val, []byte{'"'}, []byte{'\\', '"'}, 0))
				node.html.WriteString("\"")
			}
		}
	}
	node.html.WriteString("}}")
	return nil
}

// Redner
func (t *Transformater) RenderToString() string {
	return t.tree.Render()
}

func (t *Transformater) RenderBlocks() map[string]string {
	if t.blocks != nil {
		returns := map[string]string{}
		for blockId, node := range t.blocks {
			returns[blockId] = node.Render()
		}
		return returns
	}
	return nil
}

// --------------------------------------------------------------------------------

// variables
var printValueRegex, _ = regexp.Compile("^(.*){{(.*)}}$")

// func (t *Transformater) Render() string {
// 	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 	fmt.Println(t.b.String())
// 	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++=+++")
// 	return t.b.String()
// }

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
