package lock

import (
	"fmt"
	"sync"
	"testing"
)

type TBZ struct {
	L       Lock
	Counter int64
}

func (t *TBZ) DoBz() {

}

func TestTryLock(t *testing.T) {
	tt := TBZ{}
	tt.L = NewLock()
	w := sync.WaitGroup{}
	for tt.Counter < 10000 {
		w.Add(1)
		go func() {
			defer w.Done()
			if !tt.L.Lock() {
				return
			}
			tt.Counter++
			tt.L.Unlock()
		}()
	}
	w.Wait()
	fmt.Println(tt.Counter)
	//=== RUN   TestTryLock
	//10002
	//--- PASS: TestTryLock (0.01s)
	//达到锁的功能即可
}
