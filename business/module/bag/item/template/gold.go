package template

type GoldType int

const (
	NormalGold GoldType = iota + 1
	Diamond             = 2
)

type Gold struct {
	Category GoldType `json:"category"`
	*ItemBase
}
