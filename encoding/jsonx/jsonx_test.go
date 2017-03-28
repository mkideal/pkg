package jsonx

import (
	"os"
	"strings"
	"testing"
)

func ExampleDecode() {
	r := strings.NewReader(`{"a":1,"b":true,"c":[{"x":1.2},{"y":2.3}],"d":{},"e":-1,"f":+1}`)
	node, err := Read(r)
	if err != nil {
		return
	}
	Write(os.Stdout, node, WithIndent("  "))
	// Output:
	// {
	//   "a": 1,
	//   "b": true,
	//   "c": [
	//     {
	//       "x": 1.2
	//     },
	//     {
	//       "y": 2.3
	//     }
	//   ],
	//   "d": {},
	//   "e": -1,
	//   "f": +1
	// }
}

func ExampleDecodeExtraComma() {
	r := strings.NewReader(`{"a":1,"b":true,"c":[{"x":1.2},{"y":2.3},],"d":{},}`)
	node, err := Read(r, WithExtraComma())
	if err != nil {
		return
	}
	Write(os.Stdout, node, WithIndent("  "))
	Write(os.Stdout, node, WithIndent("  "), WithExtraComma())
	// Output:
	// {
	//   "a": 1,
	//   "b": true,
	//   "c": [
	//     {
	//       "x": 1.2
	//     },
	//     {
	//       "y": 2.3
	//     }
	//   ],
	//   "d": {}
	// }{
	//   "a": 1,
	//   "b": true,
	//   "c": [
	//     {
	//       "x": 1.2,
	//     },
	//     {
	//       "y": 2.3,
	//     },
	//   ],
	//   "d": {},
	// }
}

func ExampleDecodeWithUnquotedKey() {
	r := strings.NewReader(`{a:1,b:true,c:[{x:1.2},{y:2.3}],d:{}}`)
	node, err := Read(r, WithUnquotedKey())
	if err != nil {
		return
	}
	Write(os.Stdout, node, WithIndent("  "))
	Write(os.Stdout, node, WithIndent("  "), WithUnquotedKey())
	// Output:
	// {
	//   "a": 1,
	//   "b": true,
	//   "c": [
	//     {
	//       "x": 1.2
	//     },
	//     {
	//       "y": 2.3
	//     }
	//   ],
	//   "d": {}
	// }{
	//   a: 1,
	//   b: true,
	//   c: [
	//     {
	//       x: 1.2
	//     },
	//     {
	//       y: 2.3
	//     }
	//   ],
	//   d: {}
	// }
}

func ExampleDecodeWithComment() {
	r := strings.NewReader(`{
	// doc a
	"a":1, // line a
	// doc b
	"b":true, // line b

	// doc c
	"c":[
		{"x":1.2},
		// doc y
		{"y":2.3}
	], // line c
	// doc d
	"d":{}
}`)
	node, err := Read(r, WithComment())
	if err != nil {
		return
	}
	Write(os.Stdout, node, WithIndent("  "), WithComment())
	// Output:
	// {
	//   // doc a
	//   "a": 1,// line a
	//   // doc b
	//   "b": true,// line b
	//   // doc c
	//   "c": [
	//     {
	//       "x": 1.2
	//     },
	//       // doc y
	//     {
	//       "y": 2.3
	//     }
	//   ],// line c
	//   // doc d
	//   "d": {}
	// }
}

func TestParser(t *testing.T) {
	type argt struct {
		src  string
		err  string
		kind NodeKind
		opt  options
	}
	for i, ts := range []argt{
		{``, "unexpected begin of json node  at <input>:1:1", InvalidNode, options{}},
		{`%`, "unexpected begin of json node % at <input>:1:2", InvalidNode, options{}},
		{`(`, "unexpected begin of json node ( at <input>:1:2", InvalidNode, options{}},
		{`{]`, "expect a string or `}`, but got `]` at <input>:1:3", InvalidNode, options{}},
		{`//comment`, "unexpected begin of json node / at <input>:1:2", InvalidNode, options{}},
		{`/*comment*/`, "unexpected begin of json node / at <input>:1:2", InvalidNode, options{}},
		{`1`, "", IntNode, options{}},
		{`1.2`, "", FloatNode, options{}},
		{`/*comment*/1.2`, "", FloatNode, options{supportComment: true}},
		{`abc`, "", IdentNode, options{}},
		{`abc//comment`, "", IdentNode, options{supportComment: true}},
		{`'a'`, "", CharNode, options{}},
		{`''`, "illegal char literal at <input>:1:2", InvalidNode, options{}},
		{`'xxx'`, "illegal char literal at <input>:1:5", InvalidNode, options{}},
		{`""`, "", StringNode, options{}},
		{`"abcd"`, "", StringNode, options{}},
		{`'abcd"`, "illegal char literal at <input>:1:7", InvalidNode, options{}},
		{`// doc
		"abcd"`, "", StringNode, options{supportComment: true}},
		{`{"x":1}`, "", ObjectNode, options{}},
		{`{"x":1,}`, "extra comma found at <input>:1:8", InvalidNode, options{}},
		{`{"x":1,}`, "", ObjectNode, options{extraComma: true}},
		{`{"x":1,"y":{}}`, "", ObjectNode, options{}},
		{`{"x":1,"y":{]}`, "expect a string or `}`, but got `]` at <input>:1:14", InvalidNode, options{}},
		{`{"x":1,"y":{/**/}}`, "", ObjectNode, options{supportComment: true}},
		{`{"x":1,"y":{//}}`, "expect `}`, but got EOF at <input>:1:17", InvalidNode, options{supportComment: true}},
		{`[]`, "", ArrayNode, options{}},
		{`[x]`, "", ArrayNode, options{}},
		{`[x, y, z]`, "", ArrayNode, options{}},
		{`["x", "y", z]`, "", ArrayNode, options{}},
		{`[{}]`, "", ArrayNode, options{}},
		{`[1,{}]`, "", ArrayNode, options{}},
		{`[-1,{}]`, "", ArrayNode, options{}},
		{`{x:1}`, "expect a string or `}`, but got `x` at <input>:1:3", InvalidNode, options{}},
		{`{x:1}`, "", ObjectNode, options{unquotedKey: true}},
	} {
		r := strings.NewReader(ts.src)
		node, err := Read(r, ts.opt.clone)
		if err != nil && ts.err == "" {
			t.Errorf("%dth: want nil, got error %v", i, err)
			continue
		}
		if err == nil && ts.err != "" {
			t.Errorf("%dth: want error, but got nil", i)
			continue
		}
		if err != nil && ts.err != err.Error() {
			t.Errorf("%dth: want error %v, but got %v", i, ts.err, err)
			continue
		}
		if err == nil {
			if node.Kind() != ts.kind {
				t.Errorf("%dth: want node kind %v, but got %v", i, ts.kind, node.Kind())
				continue
			}
		}
	}
}
