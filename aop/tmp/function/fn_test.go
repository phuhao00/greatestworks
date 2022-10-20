package function

import "testing"

func TestFn(t *testing.T) {
	tm := Handler(KK)
	tm(1)

}
