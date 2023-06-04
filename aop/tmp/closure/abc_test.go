package closure

import (
	"fmt"
	"testing"
)

func TestNN(t *testing.T) {
	fn := func(int2 int64) {
		fmt.Println(int2)
	}
	b := []uint64{1, 2, 3, 5, 6, 7, 8, 9}
	for i := 0; i < 6; i++ {
		go fn(int64(b[i]))
	}

}

func TestNN1(t *testing.T) {
	var prints []func()
	for _, v := range []int{1, 2, 3} {
		prints = append(prints, func() { fmt.Println(v) })
	}
	for _, print := range prints {
		print()
	}
}
