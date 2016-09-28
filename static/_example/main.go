package main

import (
	"github.com/mkideal/pkg/static"
)

const (
	ErrNotFound = static.Error("not found")
)

func main() {
	e := static.Error("not found")
	e = "xxx"
	println(e)

	static.Assert(false, "assert fail: want true")
	static.Assert(ErrNotFound, "assert fail: want nil")
}
