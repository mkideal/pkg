package build

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var (
	version, branch, commit, date, time string
)

func String(name string) string {
	return fmt.Sprintf("%s %s(%s: %s) built at %s %s by %s", name, version, branch, commit, date, time, runtime.Version())
}

func Print() {
	fmt.Println(String(filepath.Base(os.Args[0])))
}

// example Makefile:
//
//	PKG=github.com/mkideal/pkg/build
//	BRANCH=$(shell git symbolic-ref --short HEAD)
//	COMMIT=$(shell git rev-parse HEAD)
//	DATE=$(shell date "+%Y/%m/%d")
//	TIME=$(shell date "+%H:%M:%S")
//
//	go build -ldflags "-X ${PKG}.branch=${BRANCH} -X ${PKG}.commit=${COMMIT} -X ${PKG}.date=${DATE} -X ${PKG}.time=${TIME}"
//
