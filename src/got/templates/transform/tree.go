package transform

import (
	"bytes"
	"fmt"
)

// Struct Tree Node
// TODO: MEM performance issue, here stores 4 copies of template files.
type Node struct {
	tagName string
	attrs   map[string][]byte
	raw     []byte
	html    bytes.Buffer

	level    int
	parent   *Node
	children []*Node
	closed   bool
}

func newNode() *Node {
	return &Node{level: 0, closed: true}
}

func (n *Node) AddChild(node *Node) {
	node.parent = n
	node.level = n.level + 1
	if n.children == nil {
		n.children = make([]*Node, 0, 2)
	}
	n.children = append(n.children, node)
}

func (n *Node) Detach() *Node {
	p := n.parent
	n.parent = nil
	if p == nil {
		return nil
	}
	for i := len(p.children) - 1; i >= 0; i-- {
		// for _, node := range p.children {
		node := p.children[i]
		// fmt.Println("find * ", node)
		if node == n {
			// fmt.Println("matched")
			p.children[i] = nil
			// p.children = append(p.children[:i], p.children[i+1:]...)
			return node
		}
	}
	return nil
}

func (n *Node) String() string {
	cn := 0
	if n.children != nil {
		cn = len(n.children)
	}
	return fmt.Sprintf("c(%v):%v", cn, n.html.String())
}

func (n *Node) Render() string {
	var html bytes.Buffer
	render(&html, n)
	return html.String()
}

func render(html *bytes.Buffer, n *Node) {
	if n == nil {
		return
	}
	html.Write(n.html.Bytes())
	if n.children != nil {
		for _, node := range n.children {
			render(html, node)
		}
	}
}

// func (n *Node) Render() string {
// 	var html bytes.Buffer
// 	render(&html, n, 0)
// 	return html.String()
// }

// func render(html *bytes.Buffer, n *Node, level int) {
// 	if n == nil {
// 		return
// 	}
// 	for i := 0; i < level; i++ {
// 		html.WriteString("  ")
// 	}
// 	// html.WriteString("[")
// 	// html.WriteString(strconv.Itoa(n.level))
// 	// if !n.closed {
// 	// 	html.WriteString("+")
// 	// }
// 	// html.WriteString("]")

// 	html.WriteString("+[")
// 	html.Write(bytes.TrimSpace(n.html.Bytes()))
// 	html.WriteString("]")

// 	// html.WriteString("  -parent:")
// 	// if n.parent == nil {
// 	// 	html.WriteString("nil")
// 	// } else {
// 	// 	html.WriteString(n.parent.tagName)
// 	// }

// 	// if n.children != nil {
// 	// 	html.WriteString("  -children:")
// 	// 	html.WriteString("<")
// 	// 	for _, nn := range n.children {
// 	// 		html.WriteString(nn.tagName)
// 	// 		html.WriteString(",")
// 	// 	}
// 	// 	html.WriteString(">")
// 	// }

// 	html.WriteString("\n")

// 	if n.children != nil {
// 		for _, node := range n.children {
// 			render(html, node, level+1)
// 		}
// 	}
// }
