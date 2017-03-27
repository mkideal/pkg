package jsonx

import (
	"bytes"
	"io"
	"text/scanner"

	"github.com/mkideal/pkg/encoding"
)

// NodeKind represents kind of json
type NodeKind int

const (
	IdentNode  NodeKind = iota // abc,true,false
	IntNode                    // 1
	FloatNode                  // 1.2
	CharNode                   // 'c'
	StringNode                 // "xyz"
	ObjectNode                 // {}
	ArrayNode                  // []
)

const (
	opLBrace = '{'
	opRBrace = '}'
	opLBrack = '['
	opRBrack = ']'
	opComma  = ','
	opColon  = ':'
)

// Node represents top-level json object
type Node interface {
	// embed encoding.Node
	encoding.Node
	// Kind returns kind of node
	Kind() NodeKind
	// Doc returns lead comments
	Doc() *encoding.CommentGroup
	// Comment returns line comments
	Comment() *encoding.CommentGroup
	// NumChild returns number of child nodes
	NumChild() int
	// ByIndex gets ith child node, key is empty if current node is not an object node
	ByIndex(i int) (key string, node Node)
	// ByKey gets child node by key, nil returned if key not found
	ByKey(key string) Node
	// Decode decodes node to ptr
	Decode(ptr interface{}) error

	// setDoc sets doc comment group
	setDoc(doc *encoding.CommentGroup)
	// setComment sets line comment group
	setComment(comment *encoding.CommentGroup)
	// output writes Node to writer
	output(prefix string, w io.Writer, opt options, lastNode bool) error
}

// Option represents a function for setting options
type Option func(*options)

type options struct {
	// prefix, indent used to pretty output json
	prefix, indent string
	// c-style comment supported if supportComment is true
	supportComment bool
	// key should be unquoted if unquotedKey is true
	unquotedKey bool
	// extra comma could be insert to end of last node of object or array if extraComma is true
	extraComma bool
}

// WithComment returns an option which sets supportComment true
func WithComment() Option {
	return func(opt *options) {
		opt.supportComment = true
	}
}

// WithPrefix returns an option which with prefix while outputing
func WithPrefix(prefix string) Option {
	return func(opt *options) {
		opt.prefix = prefix
	}
}

// WithIndent return an option which with indent while outputing
func WithIndent(indent string) Option {
	return func(opt *options) {
		opt.indent = indent
	}
}

// WithUnquotedKey returns an option which sets unquotedKey true
func WithUnquotedKey() Option {
	return func(opt *options) {
		opt.unquotedKey = true
	}
}

// WithExtraComma returns an option which sets extraComma true
func WithExtraComma() Option {
	return func(opt *options) {
		opt.extraComma = true
	}
}

func applyOptions(opts []Option) options {
	opt := options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Read reads a json node from reader r
func Read(r io.Reader, opts ...Option) (Node, error) {
	opt := applyOptions(opts)
	s := new(scanner.Scanner)
	s = s.Init(r)
	s.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanChars | scanner.ScanStrings
	if opt.supportComment {
		s.Mode |= scanner.ScanComments
	}
	p := new(parser)
	p.init(s, opt)
	return p.parseNode()
}

// Write write a json node to writer w
func Write(w io.Writer, node Node, opts ...Option) error {
	return node.output("", w, applyOptions(opts), true)
}

// Unmarshal unmarshals data to pointer ptr
func Unmarshal(data []byte, ptr interface{}, opts ...Option) error {
	node, err := Read(bytes.NewBuffer(data), opts...)
	if err == nil {
		err = node.Decode(ptr)
	}
	return err
}

//TODO: implements Marshal function: func Marshal(i interface{}, opts ...Option) ([]byte, error)
