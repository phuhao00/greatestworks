package colors

import (
	"os"
	"testing"
)

func TestTabularizer(t *testing.T) {
	red := Color256(124)
	blue := Color256(27)
	green := Color256(40)

	title := []Text{
		{{S: "CATS", Bold: true}},
		{{S: "felis catus", Color: Color256(254)}},
	}
	tab := NewTabularizer(os.Stdout, title, NoDim)
	defer tab.Flush()
	tab.Row("NAME", "AGE", "COLOR")
	tab.Row("belle", "1y", Atom{S: "tortie", Color: red})
	tab.Row("sidney", "2y", Atom{S: "calico", Color: blue})
	tab.Row("dakota", "8m", Atom{S: "tuxedo", Color: green, Underline: true})
}
