package jsonx

import (
	"os"
	"strings"
)

func ExampleDecode() {
	r := strings.NewReader(`{"a":1,"b":true,"c":[{"x":1.2},{"y":2.3}],"d":{}}`)
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
	//   "d": {}
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
