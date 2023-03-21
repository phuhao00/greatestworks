package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Router struct {
	real *mux.Router
}

func (r *Router) Register(path string, fn func(http.ResponseWriter, *http.Request)) {
	r.real.HandleFunc(path, fn)
}

func (r *Router) Run() {

}

func (r *Router) Init() {
	r.real = mux.NewRouter()
	r.Register("character/set_vip", SetVip)
}

func SetVip(http.ResponseWriter, *http.Request) {
	//todo
}
