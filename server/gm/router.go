package main

import (
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/proto"
	"greatestworks/server/gm/user"
	"greatestworks/server/gm/vip"
	"net/http"
)

type Router struct {
	real      *mux.Router
	toGateWay chan proto.Message
}

func (r *Router) AddHandler(path string, fn func(http.ResponseWriter, *http.Request)) {
	r.real.HandleFunc(path, fn)
}

func (r *Router) Run() {
	r.Init()

}

func (r *Router) Init() {
	r.real = mux.NewRouter()
	r.AddHandler("user/register", user.Register)
	r.AddHandler("user/login", user.Login)
	r.AddHandler("character/set_vip", vip.SetVip)
}
