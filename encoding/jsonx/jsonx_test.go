package jsonx

import (
	"os"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	r := strings.NewReader(`{"a":1,"b":true,"c":[{"x":1.2},{"y":2.3}],"d":{}}`)
	node, err := Read(r, WithExtraComma())
	if err != nil {
		t.Fatalf("error: %v", err)
		return
	}
	Write(os.Stdout, node, WithIndent("    "))
}

func TestDecodeWithUnquotedKey(t *testing.T) {
	r := strings.NewReader(`{a:1,b:true,c:[{x:1.2},{y:2.3}],d:{}}`)
	node, err := Read(r, WithUnquotedKey())
	if err != nil {
		t.Fatalf("error: %v", err)
		return
	}
	Write(os.Stdout, node, WithIndent("    "))
}

func TestDecodeWithComment(t *testing.T) {
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
		t.Fatalf("error: %v", err)
		return
	}
	Write(os.Stdout, node, WithIndent("    "), WithComment())
}
