package goid

import (
	"fmt"
	"testing"
)

func TestGetGoId(t *testing.T) {
	go func() {
		fmt.Println("协程1：", GetGoId())

	}()

	go func() {
		fmt.Println("协程2：", GetGoId())

	}()

	go func() {
		fmt.Println("协程3：", GetGoId())

	}()

	//select {}
}

func TestGetGoIdWithReflect(t *testing.T) {
	go func() {
		fmt.Println("协程1：", GetGoIdWithReflect())

	}()

	go func() {
		fmt.Println("协程2：", GetGoIdWithReflect())

	}()

	go func() {
		fmt.Println("协程3：", GetGoIdWithReflect())

	}()

}
