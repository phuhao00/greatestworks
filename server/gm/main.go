package main

import "github.com/gorilla/mux"

func main() {
	h := Router{
		real: mux.NewRouter(),
	}
	h.Run()
}
