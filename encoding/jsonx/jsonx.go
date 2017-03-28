package jsonx

import (
	"io"
	"text/scanner"
)

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

func (opt options) clone(dst *options) {
	dst.prefix = opt.prefix
	dst.indent = opt.indent
	dst.supportComment = opt.supportComment
	dst.unquotedKey = opt.unquotedKey
	dst.extraComma = opt.extraComma
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
	if err := p.init(s, opt); err != nil {
		return nil, err
	}
	return p.parseNode()
}

// Write writes a json node to writer w
func Write(w io.Writer, node Node, opts ...Option) error {
	return node.output("", w, applyOptions(opts), true, true)
}
