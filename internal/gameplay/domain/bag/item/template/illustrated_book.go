package template

type IllustratedBookType int

const (
	Fish   IllustratedBookType = iota + 1
	Flower                     = 2
)

type IllustratedBook struct {
	*ItemBase
	Category IllustratedBookType
}
