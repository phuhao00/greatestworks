package lua

import (
	"github.com/yuin/gopher-lua"
)

func GetLua() {
	L := lua.NewState()
	defer L.Close()
	if err := L.DoFile("hello.lua"); err != nil {
		panic(err)
	}
}
