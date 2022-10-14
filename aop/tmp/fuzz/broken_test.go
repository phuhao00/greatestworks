package fuzz

import (
	"fmt"
	"testing"
)

func FuzzBrokenMethod(f *testing.F) {
	f.Fuzz(func(t *testing.T, Data string) {
		BrokenMethod(Data)
	})
}

func FuzzMod(f *testing.F) {
	f.Fuzz(func(t *testing.T, a, b int) {
		fmt.Println(a / b)
	})
}

func FuzzReverse(f *testing.F) {
	f.Fuzz(func(t *testing.T, a string) {
		Reverse(a)
	})
}

// go test -fuzztime 10s -fuzz  FuzzMod

//https://go.dev/security/fuzz/
