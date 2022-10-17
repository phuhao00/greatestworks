package main

import (
	"greatestworks/aop/pprof"
	"greatestworks/aop/pprof/web"
)

//
func main() {
	h := pprof.Handler{
		Router: web.NewHttpRouter(),
	}
	h.RegisterProfiler()

}
