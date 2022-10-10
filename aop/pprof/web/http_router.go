package web

import (
	"github.com/gorilla/mux"
	"net/http"
	"sync/atomic"
)

// HttpRouter is a Router implementation for the Gorilla web toolkit's `mux.Router`.
type HttpRouter struct {
	mux      *mux.Router
	pageView int64 // requests process by multi goroutine，need atomic
}

func NewHttpRouter() *HttpRouter {
	return &HttpRouter{mux: mux.NewRouter().StrictSlash(true), pageView: 0}
}

// Handle will call the Gorilla web toolkit's Handle().Method() methods.
func (g *HttpRouter) Handle(method, path string, h http.Handler) {
	g.mux.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// copy the route params into a shared location
		// duplicating memory, but allowing Gizmo to be more flexible with
		// router implementations.
		atomic.AddInt64(&g.pageView, 1)
		SetRouteVars(r, mux.Vars(r))
		h.ServeHTTP(w, r)
	})).Methods(method)
}

func (g *HttpRouter) GetPageView() int64 {
	return atomic.LoadInt64(&g.pageView)
}

// HandleFunc will call the Gorilla web toolkit's HandleFunc().Method() methods.
func (g *HttpRouter) HandleFunc(method, path string, h func(http.ResponseWriter, *http.Request)) {
	g.Handle(method, path, http.HandlerFunc(h))
}

// HandleStaticFile will call the Gorilla web toolkit's HandleFunc().Method() methods.
func (g *HttpRouter) HandleStaticFile(method, path, dir string) {
	g.mux.PathPrefix(path).Handler(http.StripPrefix(path, http.FileServer(http.Dir(dir))))
}

// SetNotFoundHandler will set the Gorilla mux.Router.NotFoundHandler.
func (g *HttpRouter) SetNotFoundHandler(h http.Handler) {
	g.mux.NotFoundHandler = h
}

// ServeHTTP will call Gorilla mux.Router.ServerHTTP directly.
func (g *HttpRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}
