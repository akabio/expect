package expect

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"gitlab.com/akabio/gopath"
)

// location decorator, this files location is used to report errors
// for this reason the name is x.go so it's less intrusive then the
// real location

var warned = false

type locationDeco struct {
	t Test
}

// extra space so line numbers can be adapted
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//

func (lt *locationDeco) Fatalf(f string, i ...interface{}) {
	fnc, file, line, ok := getLocation()
	if ok {
		loc := fmt.Sprintf("%v:%v %v\n", file, line, fnc)
		lt.t.Fatalf("%v"+f, prep(loc, i)...)
		return
	}
	lt.t.Fatalf(f, i...)
}

func (lt *locationDeco) Errorf(f string, i ...interface{}) {
	fnc, file, line, ok := getLocation()
	if ok {
		loc := fmt.Sprintf("%v:%v %v\n", file, line, fnc)
		lt.t.Errorf("%v"+f, prep(loc, i)...)
		return
	}
	lt.t.Errorf(f, i...)
}

func (lt *locationDeco) Error(p ...interface{}) {
	fnc, file, line, ok := getLocation()
	if ok {
		loc := fmt.Sprintf("%v:%v %v\n", file, line, fnc)
		lt.t.Error(prep(loc, p)...)
		return
	}
	lt.t.Error(p...)
}

func prep(m string, i []interface{}) []interface{} {
	return append(append([]interface{}{}, m), i...)
}

func shortFunc(name string) string {
	// skip all except package name
	lastSlash := strings.LastIndex(name, "/")
	if lastSlash != -1 {
		name = name[lastSlash+1:]
	}
	// skip package name
	firstDot := strings.Index(name, ".")
	if firstDot != -1 {
		name = name[firstDot+1:]
	}
	return name
}

func getLocation() (string, string, int, bool) {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return "", "", 0, false
	}

	file, err := gopath.RelativePath(file)
	if err != nil && !warned {
		log.Println("no go.mod found, can't calculate relative path")
		warned = true
	}

	pc, _, _, ok := runtime.Caller(3)
	if !ok {
		return "", "", 0, false
	}

	fnc := shortFunc(runtime.FuncForPC(pc).Name())
	return fnc, file, line, true
}
