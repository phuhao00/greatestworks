package plant

type Config struct {
	Id   uint32 `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type Status uint16

const (
	Growing Status = iota + 1
	Maturation
)
