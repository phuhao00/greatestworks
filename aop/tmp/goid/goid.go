package goid

import (
	"bytes"
	"reflect"
	"runtime"
	"strconv"
)

func GetGoId() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	//fmt.Println("stack", string(b))
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func getG() interface{}

func GetGoIdWithReflect() int64 {
	g := getG()
	goId := reflect.ValueOf(g).FieldByName("goid").Int()
	return goId
}
