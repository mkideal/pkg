package jsonx

import (
	"fmt"
	"io"
	"strconv"
	"text/scanner"

	"github.com/mkideal/pkg/encoding"
)

func outputDoc(prefix string, w io.Writer, doc *encoding.CommentGroup) error {
	if doc == nil {
		return nil
	}
	var err error
	for i, line := range doc.List {
		if i > 0 {
			fmt.Fprint(w, "\n"+prefix+line.Text)
		} else {
			fmt.Fprint(w, prefix+line.Text)
		}
		if err != nil {
			return err
		}
	}
	return err
}

func outputNext(prefix string, w io.Writer, opt options) error {
	if opt.indent == "" {
		return nil
	}
	_, err := fmt.Fprint(w, "\n"+prefix+opt.indent)
	return err
}

func outputNodeTail(w io.Writer, n Node, topNode, lastNode bool, opt options) error {
	if (opt.extraComma || !lastNode) && !topNode {
		if _, err := fmt.Fprint(w, ","); err != nil {
			return err
		}
	}
	if opt.supportComment && n.Comment() != nil {
		if _, err := fmt.Fprint(w, n.Comment().Text()); err != nil {
			return err
		}
	}
	return nil
}

// kv represents a key-value pair
type kv struct {
	key   string
	value Node
}

// nodebase represents base of any json node
type nodebase struct {
	pos     scanner.Position
	doc     *encoding.CommentGroup
	comment *encoding.CommentGroup
}

func (n nodebase) Pos() scanner.Position           { return n.pos }
func (n nodebase) Doc() *encoding.CommentGroup     { return n.doc }
func (n nodebase) Comment() *encoding.CommentGroup { return n.comment }
func (n nodebase) NumChild() int                   { return 0 }
func (n nodebase) ByIndex(i int) (string, Node)    { return "", nil }
func (n nodebase) ByKey(key string) Node           { return nil }

func (n *nodebase) setDoc(doc *encoding.CommentGroup)         { n.doc = doc }
func (n *nodebase) setComment(comment *encoding.CommentGroup) { n.comment = comment }

// objectNode represents object node
type objectNode struct {
	nodebase
	children []kv
	indexMap map[string]int
}

func newObjectNode() *objectNode {
	return new(objectNode)
}

func (n *objectNode) addChild(key string, value Node) {
	if n.indexMap == nil {
		n.indexMap = make(map[string]int)
	}
	index, ok := n.indexMap[key]
	if !ok {
		n.children = append(n.children, kv{key, value})
	} else {
		n.children[index].value = value
	}
}

func (n objectNode) Interface() interface{} {
	m := make(map[string]interface{})
	for _, kv := range n.children {
		m[kv.key] = kv.value.Interface()
	}
	return m
}

func (n objectNode) Kind() NodeKind               { return ObjectNode }
func (n objectNode) NumChild() int                { return len(n.children) }
func (n objectNode) ByIndex(i int) (string, Node) { return n.children[i].key, n.children[i].value }
func (n objectNode) ByKey(key string) Node {
	if n.indexMap == nil {
		return nil
	}
	index, ok := n.indexMap[key]
	if !ok {
		return nil
	}
	return n.children[index].value
}

func (n *objectNode) output(prefix string, w io.Writer, opt options, topNode, lastNode bool) error {
	writeComment := opt.indent != "" && opt.supportComment
	if _, err := fmt.Fprint(w, "{"); err != nil {
		return err
	}
	numChild := len(n.children)
	for i, child := range n.children {
		if err := outputNext(prefix, w, opt); err != nil {
			return err
		}
		doc := child.value.Doc()
		if writeComment && doc != nil {
			if err := outputDoc(prefix, w, doc); err != nil {
				return err
			}
			if err := outputNext(prefix, w, opt); err != nil {
				return err
			}
		}
		key := child.key
		// try quote key string with "
		if len(key) > 0 && key[0] != '"' && !opt.unquotedKey {
			key = strconv.Quote(key)
		}
		if _, err := fmt.Fprint(w, key+":"); err != nil {
			return err
		}
		if opt.indent != "" {
			// insert a space between key-value pair
			if _, err := fmt.Fprint(w, " "); err != nil {
				return err
			}
		}
		if err := child.value.output(prefix+opt.indent, w, opt, false, i+1 == numChild); err != nil {
			return err
		}
	}
	next := "}"
	if len(n.children) > 0 {
		next = prefix + next
		if opt.indent != "" {
			next = "\n" + next
		}
	}
	if _, err := fmt.Fprint(w, next); err != nil {
		return err
	}
	return outputNodeTail(w, n, topNode, lastNode, opt)
}

// arrayNode represents array node
type arrayNode struct {
	nodebase
	children []Node
}

func newArrayNode() *arrayNode {
	return new(arrayNode)
}

func (n *arrayNode) addChild(value Node) {
	n.children = append(n.children, value)
}

func (n arrayNode) Interface() interface{} {
	size := len(n.children)
	if size == 0 {
		return []interface{}{}
	}
	s := make([]interface{}, 0, size)
	for _, child := range n.children {
		s = append(s, child.Interface())
	}
	return s
}

func (n arrayNode) Kind() NodeKind               { return ArrayNode }
func (n arrayNode) NumChild() int                { return len(n.children) }
func (n arrayNode) ByIndex(i int) (string, Node) { return "", n.children[i] }
func (n arrayNode) ByKey(key string) Node        { return nil }

func (n *arrayNode) output(prefix string, w io.Writer, opt options, topNode, lastNode bool) error {
	writeComment := opt.indent != "" && opt.supportComment
	if _, err := fmt.Fprint(w, "["); err != nil {
		return err
	}
	numChild := len(n.children)
	for i, child := range n.children {
		if err := outputNext(prefix, w, opt); err != nil {
			return err
		}
		doc := child.Doc()
		if writeComment && doc != nil {
			if err := outputDoc(prefix, w, doc); err != nil {
				return err
			}
			if err := outputNext(prefix, w, opt); err != nil {
				return err
			}
		}
		if err := child.output(prefix+opt.indent, w, opt, false, i+1 == numChild); err != nil {
			return err
		}
	}
	next := "]"
	if len(n.children) > 0 {
		next = prefix + next
		if opt.indent != "" {
			next = "\n" + next
		}
	}
	if _, err := fmt.Fprint(w, next); err != nil {
		return err
	}
	return outputNodeTail(w, n, topNode, lastNode, opt)
}

// basicNode represents a basic node, e.g. char,string,ident,float,int
type basicNode struct {
	nodebase
	kind  NodeKind
	value string
}

func newBasicNode(pos scanner.Position, tok rune, value string) (*basicNode, error) {
	n := &basicNode{
		nodebase: nodebase{
			pos: pos,
		},
		value: value,
	}
	switch tok {
	case scanner.Char:
		n.kind = CharNode
	case scanner.String:
		n.kind = StringNode
	case scanner.Float:
		n.kind = FloatNode
	case scanner.Int:
		n.kind = IntNode
	case scanner.Ident:
		n.kind = IdentNode
	default:
		return nil, fmt.Errorf("unexpected begin of json node %v at %v", value, pos)
	}
	return n, nil
}

func (n basicNode) Interface() interface{} {
	switch n.kind {
	case CharNode:
		value, _, _, _ := strconv.UnquoteChar(n.value, '\'')
		return value
	case StringNode:
		value, _ := strconv.Unquote(n.value)
		return value
	case FloatNode:
		value, _ := strconv.ParseFloat(n.value, 64)
		return value
	case IntNode:
		value, _ := strconv.ParseInt(n.value, 0, 64)
		return value
	case IdentNode:
		return n.value
	default:
		return nil
	}
}

func (n *basicNode) Kind() NodeKind { return n.kind }

func (n *basicNode) output(prefix string, w io.Writer, opt options, topNode, lastNode bool) error {
	if _, err := fmt.Fprint(w, n.value); err != nil {
		return err
	}
	return outputNodeTail(w, n, topNode, lastNode, opt)
}
