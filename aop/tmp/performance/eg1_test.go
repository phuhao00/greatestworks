package performance

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestZeroStruct(t *testing.T) {
	var a int
	var b string
	var e struct{}
	fmt.Println(unsafe.Sizeof(a)) // 4
	fmt.Println(unsafe.Sizeof(b)) // 8
	fmt.Println(unsafe.Sizeof(e)) // 0

}
