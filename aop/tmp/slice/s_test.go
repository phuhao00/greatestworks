package slice

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSlice(t *testing.T) {
	//S()
	v := [3]int{}
	fmt.Println(reflect.TypeOf(v).Kind())
}
