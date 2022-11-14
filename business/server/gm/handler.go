package main

import (
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	real *mux.Router
}

func (h *Handler) Register(path string, fn func(http.ResponseWriter, *http.Request)) {
	h.real.HandleFunc(path, fn)
}

func (h *Handler) Run() {

}
