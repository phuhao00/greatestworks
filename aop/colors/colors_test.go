package colors

import (
	"fmt"
	"testing"
)

// TestPrintVisibleColors prints out the colors in visible. If you run this
// test with -v, you can see what all the colors look like.
func TestPrintVisibleColors(t *testing.T) {
	s := ""
	for i, code := range visible {
		s += fmt.Sprintf("%s%4d%s ", Color256(code), code, Reset)
		if (i+1)%10 == 0 {
			t.Log(s)
			s = ""
		}
	}
	if s != "" {
		t.Log(s)
	}
}
