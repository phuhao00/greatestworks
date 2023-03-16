package iter

import (
	"fmt"
	"sync"
	"testing"
)

func TestName(t *testing.T) {
	ForRangeAsyncFix2()
}

func ForRangeAsyncFix2() {
	type Person struct {
		Name string
	}
	group := []*Person{{"li"}, {"zhao"}}

	var wg sync.WaitGroup

	for _, p := range group {
		p := p
		wg.Add(1)

		go func() {
			defer wg.Done()
			fmt.Println("name =", p.Name)
		}()
	}
	wg.Wait()
}

// name = zhao
// name = li
