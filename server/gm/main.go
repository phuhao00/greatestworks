package main

import "github.com/gorilla/mux"

func main() {
	h := Handler{
		real: mux.NewRouter(),
	}
	h.Run()
}
