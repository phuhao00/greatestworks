package pprof

import (
	"greatestworks/aop/pprof/web"
	"net/http/pprof"
)

type Handler struct {
	Router *web.HttpRouter
}

func (hs *Handler) RegisterProfiler() {
	hs.Router.HandleFunc("GET", "/debug/pprof/", pprof.Index)
	hs.Router.HandleFunc("GET", "/debug/pprof/cmdline", pprof.Cmdline)
	hs.Router.HandleFunc("GET", "/debug/pprof/profile", pprof.Profile)
	hs.Router.HandleFunc("GET", "/debug/pprof/symbol", pprof.Symbol)
	hs.Router.HandleFunc("GET", "/debug/pprof/trace", pprof.Trace)

	// 研发期间使用，稳定运营阶段，优先考虑屏蔽
	// Manually add support for paths linked to by index page at /debug/pprof/
	hs.Router.Handle("GET", "/debug/pprof/goroutine", pprof.Handler("goroutine"))
	hs.Router.Handle("GET", "/debug/pprof/heap", pprof.Handler("heap"))
	hs.Router.Handle("GET", "/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	hs.Router.Handle("GET", "/debug/pprof/block", pprof.Handler("block"))
}
