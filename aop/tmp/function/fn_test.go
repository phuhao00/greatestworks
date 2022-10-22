package function

import "testing"

func TestFn(t *testing.T) {
	tm := Handler(KK)
	tm(1)

}

func TestAbc(t *testing.T) {
	var a *ABC
	a.DO()

}