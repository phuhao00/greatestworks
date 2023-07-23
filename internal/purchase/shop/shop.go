package shop

type Category int

const (
	CategoryNormal Category = iota + 1
	CategoryMystery
)

type Shop interface {
	Refresh()
}
