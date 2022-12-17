package plant

type Config struct {
	Id   uint32
	Name string
	Desc string
}

type Status uint16

const (
	Growing Status = iota + 1
	Maturation
)
