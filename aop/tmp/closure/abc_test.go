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
