/**
  Time-stamp: <[transform.go] Elivoa @ Sunday, 2013-07-28 18:16:17>
*/

package transform

import (
	"bytes"
	"code.google.com/p/go-html-transform/h5"
	"code.google.com/p/go.net/html"
	"fmt"
	"io"
	"strings"
)

// ---- Transform template ------------------------------------------
type Transformater struct {
	tree        *h5.Tree
	deleteQueue []*html.Node
}

func NewTransformer() *Transformater {
	return &Transformater{
		deleteQueue: []*html.Node{},
	}
}

// Transform tempalte fiels. functions:
// translate <t:some_component ... /> into {{t_some_component ...}}
//
func (t *Transformater) Parse(reader io.Reader) error {
	t.tree, _ = h5.New(reader)

	// 1. read import node.
	t.tree.Walk(func(n *html.Node) {
		text := strings.TrimSpace(n.Data)
		fmt.Printf("---- %v ---- %v\n", n.Type, text)

		if n.Type == html.ElementNode {
			if strings.HasPrefix(n.Data, "t:") {
				cmd := n.Data[2:]
				switch cmd {
				case "import":
					t.replaceNode(n, nil) // remove node
				default:
					// components
					componentText := transformComponent(cmd, n)
					t.replaceNode(n, UnescapedText(componentText))
				}
			}
		}
		if n.Type == html.TextNode {
			if text != "" {
				// fmt.Printf(" %v: %v  Data: %v\n", n.DataAtom, n.Namespace, text)
				// print node
				// str := RenderNodesToString([]*html.Node{n})
				// fmt.Println("---- ", str)
			}
		}
	})

	t.deleteQueueNodes() // delete node here.

	return nil
}

func (t *Transformater) replaceNode(node *html.Node, newNode *html.Node) {
	if newNode != nil {
		node.Parent.InsertBefore(newNode, node)
	}
	// node.Parent.RemoveChild(node) // remove later
	t.deleteQueue = append(t.deleteQueue, node)
}

func (t *Transformater) deleteQueueNodes() {
	for _, node := range t.deleteQueue {
		node.Parent.RemoveChild(node)
	}
}

// ---- Transformation ------------------------------------------

// replace node, return node to remove.

// name is component lookup key
func transformComponent(name string, n *html.Node) string {
	fmt.Println("---- Transform Componnet [", name, "] ------------------------")

	var t bytes.Buffer
	t.WriteString("{{t_") // prefix
	t.WriteString(name)   // component lookup key
	t.WriteString(" $ ")  // the first parameter, context

	for _, attr := range n.Attr {
		// key:
		t.WriteString(" \"")
		t.WriteString(strings.ToUpper(attr.Key[0:1]))
		t.WriteString(attr.Key[1:])
		t.WriteString("\" ")

		// Value transform: for name="_some_value_", we transform it into:
		//   ~ before ~           ~ after ~             ~ note ~
		//   ".Name"              .Name                // start form . or $
		//   "literal:....."      "...."               // literal prefix
		//	 "abcd"               "abcd"               // auto detect plan text
		//	 ".Name+'_'+.Id"      (print .Name '_' .Id)// special
		//
		//   TODO support more prefix...
		//
		// if value starts from . or $ , treate this as property. others as string
		if strings.HasPrefix(attr.Val, ".") || strings.HasPrefix(attr.Val, "$") ||
			strings.HasPrefix(attr.Val, "(") {
			t.WriteString(attr.Val)
		} else if strings.HasPrefix(attr.Val, "print ") {
			t.WriteString("(")
			t.WriteString(attr.Val[5:])
			t.WriteString(")")
		} else if strings.HasPrefix(attr.Val, "literal:") {
			t.WriteString(" `")
			t.WriteString(attr.Val[8:])
			t.WriteString("`")
		} else { // default
			t.WriteString(attr.Val)
			// // if no space, this is literal.
			// if strings.Contains(attr.Val, " ") {
			// 	t.WriteString(attr.Val)
			// } else {
			// 	t.WriteString(" `")
			// 	t.WriteString(attr.Val)
			// 	t.WriteString("`")
			// }
		}
	}
	t.WriteString("}}")
	return t.String()
}

func writeFunctionalValue(t bytes.Buffer, value string) {
	fmt.Println(">>>> ", value)
	t.WriteString(value)
}

func writePrintValue(t bytes.Buffer, value string) {
	t.WriteString("(")
	t.WriteString(value)
	t.WriteString(")")
}

func writeLiteralValue(t bytes.Buffer, value string) {
	t.WriteString(" `")
	t.WriteString(value)
	t.WriteString("`")

}

// ---- Render Nodes to html ------------------------------------------

func (t *Transformater) RenderToString() string {
	return t.RenderNodesToString([]*html.Node{t.tree.Top()})
}

func (t *Transformater) RenderNodesToString(ns []*html.Node) string {
	buf := bytes.NewBufferString("")
	RenderNodes(buf, ns)
	return string(buf.Bytes())
}

func RenderNodes(w io.Writer, ns []*html.Node) error {
	for _, n := range ns {
		err := Render(w, n)
		if err != nil {
			return err
		}
	}
	return nil
}

// create unescaped text node
func UnescapedText(str string) *html.Node {
	return &html.Node{
		Data: str,
		Type: html.TextNode,
	}
}
