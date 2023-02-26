package track

import (
	"fmt"
	"time"
)

func GetExecuteTimeWrapFunc(f func()) func() {
	return func() {
		begin := time.Now()
		defer func() {
			end := time.Now()
			fmt.Println("execute time is", end.Sub(begin).Nanoseconds())
		}()
		f()
	}
}
