package annotation

import (
	"fmt"
	"testing"
)

type RegisterMock struct {
	ModuleName string
	Cb         func() int
}

var (
	regMock    map[string]*RegisterMock
	regFnMock  RegisterNormalMock
	eachFnMock ForEachNormalMock
)

type RegisterNormalMock func(moduleName string, fn func() int)

type ForEachNormalMock func(f int) int

func TestReg(t *testing.T) {
	regMock = make(map[string]*RegisterMock)
	regFnMock = func(moduleName string, fn func() int) {
		regMock[moduleName] = &RegisterMock{
			ModuleName: moduleName,
			Cb:         fn,
		}
	}
	regFnMock("装备", func() int {
		return 1
	})
	regFnMock("宠物", func() int {
		return 2
	})
	eachFnMock = func(f int) int {
		for _, mock := range regMock {
			f += mock.Cb()
			fmt.Println(mock.ModuleName)
		}
		return f
	}
	var combat int
	combat = eachFnMock(combat)
	fmt.Println("combat:", combat)
}
