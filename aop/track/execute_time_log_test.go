package track

import (
	"fmt"
	"testing"
)

func TestExecuteTimeLog(t *testing.T) {
	GetExecuteTimeWrapFunc(func() {
		for i := 0; i < 50000; i++ {
			fmt.Println(i)
		}
	})()
}
